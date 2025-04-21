package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/candrap89/droneTravelingDistance/generated"
	"github.com/candrap89/droneTravelingDistance/handler"
	"github.com/candrap89/droneTravelingDistance/kafka"
	"github.com/candrap89/droneTravelingDistance/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"localhost:7000",
			"localhost:7001",
			"localhost:7002",
		},
		Password: "",
	})
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Use(ApiKeyMiddleware(redisClient))
	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := "postgresql://postgres:andromeda@localhost:5432/postgres?sslmode=disable"
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	opts := handler.NewServerOptions{
		Repository: repo,
	}
	// Initialize Kafka consumer
	kafkaConsumer := kafka.NewConsumerHandler(repo)
	go kafkaConsumer.StartNewProductConsumer()

	return handler.NewServer(opts)
}

func ApiKeyMiddleware(redisClient *redis.ClusterClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.Background()

			apiKey := c.Request().Header.Get("X-API-Key")
			if apiKey == "" {
				return next(c)
			}
			fmt.Println("path:", c.Path())
			path := strings.Trim(c.Path(), "/")
			keyPath := strings.ReplaceAll(path, "/", "_")
			keyPath = strings.ReplaceAll(keyPath, ":", "")
			method := strings.ToLower(c.Request().Method)
			redisKey := "api:" + method + ":" + keyPath // e.g., "api:get:estate_id_tree"
			fmt.Println("Redis Key:", redisKey)

			// 🔌 Use injected redisClient here
			storedKey, err := redisClient.Get(ctx, redisKey).Result()
			if err == redis.Nil {
				return next(c)
				//return echo.NewHTTPError(http.StatusForbidden, "API key config not found for this route")
			} else if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Redis error: "+err.Error())
			}

			if storedKey != apiKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API key")
			}

			// Rate Limiting
			// Increment the request count for the API key
			// For example, you can use Redis to track the number of requests made with this API key
			// and check if it exceeds a certain limit within a time window.
			// Example: Increment the request count in Redis
			count, err := redisClient.Incr(ctx, "rate_limit:"+apiKey).Result()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Redis error: "+err.Error())
			}
			fmt.Println("Request count for API key:", apiKey, "is", count)
			// Set an expiration time for the key (e.g., 1 hour)
			// This is a placeholder; adjust the expiration time as needed
			_, err = redisClient.Expire(ctx, "rate_limit:"+apiKey, 60*time.Second).Result()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Redis error: "+err.Error())
			}
			// Check if the rate limit is exceeded
			// For example, you can set a limit of 100 requests per hour
			// This is a placeholder; adjust the limit as needed
			limit := 100
			if count > int64(limit) {
				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			}
			// Add hit count to context
			c.Set("hit_count", count)
			fmt.Println("Hit count: ", c.Get("hit_count"))

			// Call the next handler
			return next(c)
		}
	}
}
