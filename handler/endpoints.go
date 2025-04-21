package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/candrap89/droneTravelingDistance/generated"
	"github.com/candrap89/droneTravelingDistance/kafka"
	"github.com/candrap89/droneTravelingDistance/repository"
	"github.com/candrap89/droneTravelingDistance/utils"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) GetHello(ctx echo.Context, params generated.GetHelloParams) error {
	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

type ApiResponse struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	Data       []struct {
		City          string `json:"city"`
		Name          string `json:"name"`
		EstimatedCost int32  `json:"estimated_cost"`
		UserRating    struct {
			AverageRating float64 `json:"average_rating"`
			Votes         int32   `json:"votes"`
		} `json:"user_rating"`
	} `json:"data"`
}

func (s *Server) GetVoteCount(c echo.Context, params generated.GetVoteCountParams) error {
	url := fmt.Sprintf("https://jsonmock.hackerrank.com/api/food_outlets?city=%s&estimated_cost=%d", params.CityName, params.EstimatedCost)

	// Set a timeout for the HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retry logic: 3 retries, 2s delay, exponential backoff
	resp, err := utils.DoWithRetry(ctx, "GET", url, 3, 2*time.Second, true)
	if err != nil {
		fmt.Println("GET request failed after retries:", err)
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Failed to fetch data from external API",
		})
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Error reading response body",
		})
	}

	// Parse the JSON response
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Error parsing JSON",
		})
	}

	// Check if any data matches the criteria
	if len(apiResponse.Data) == 0 {
		fmt.Println("No matching restaurant found.")
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "No matching restaurant found",
		})
	}

	var voteCount int32
	for _, data := range apiResponse.Data {
		voteCount += data.UserRating.Votes
	}

	//publish kafka message
	kafka.SendNewProductMessage(int(voteCount), params.CityName)

	// Return the vote count of the first matching restaurant
	return c.JSON(http.StatusOK, generated.CityVoteResponse{
		City:      params.CityName,
		VoteCount: int(voteCount),
	})
}

// PostEstate is a placeholder implementation to satisfy the generated.ServerInterface.
func (s *Server) PostEstate(ctx echo.Context) error {
	var input generated.Estate

	// Bind JSON input to generated.Estate struct
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid request payload",
		})
	}
	// Basic validation
	if input.Width <= 0 || input.Length <= 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Width and length must be greater than 0",
		})
	}
	fmt.Println("Width:", input.Width)
	fmt.Println("Length:", input.Length)
	// Call repository to create estate
	output, err := s.Repository.CreateEstate(ctx.Request().Context(), repository.CreateEstateInput{
		Width:  int(input.Width),
		Length: int(input.Length),
	})
	fmt.Println("Output:", err)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Failed to create estate",
		})
	}
	// Return 201 with ID
	return ctx.JSON(http.StatusOK, generated.CreateEstateResponse{
		Id: &output.ID,
	})
}

func (s *Server) PostEstateIdTree(ctx echo.Context, id types.UUID) error {
	var input generated.CreateTreeRequest
	// Bind JSON input to generated.Tree struct
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid request payload",
		})
	}
	// Basic validation
	if input.X <= 0 || input.Y <= 0 || input.Height <= 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "X and Y coordinates and Height must be greater than or equal to 0",
		})
	}
	// validate if the tree is out of bounds
	estate, err := s.Repository.GetDetailEstate(ctx.Request().Context(), repository.GetDetailEstateInput{
		ID: id,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Failed to get estate details",
		})
	}
	if input.X > estate.Length || input.Y > estate.Width {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: fmt.Sprintf("Tree coordinates (%d, %d) are out of bounds for estate (%d, %d)", input.X, input.Y, estate.Width, estate.Length),
		})
	}
	// Check if the tree already exists
	for _, tree := range estate.Trees {
		if tree.X == int(input.X) && tree.Y == int(input.Y) {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
				Message: fmt.Sprintf("Tree already exists at coordinates (%d, %d)", input.X, input.Y),
			})
		}
	}

	// Call repository to create tree
	output, err := s.Repository.CreateTree(ctx.Request().Context(), repository.CreateTreeInput{
		EstateID: id,
		X:        int(input.X),
		Y:        int(input.Y),
		Height:   int(input.Height),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Failed to create tree",
		})
	}
	// Return 201 with ID
	return ctx.JSON(http.StatusOK, generated.CreateTreeResponse{
		Id: &output.ID,
	})
}

func (s *Server) GetEstateIdStats(ctx echo.Context, id types.UUID) error {
	// Call repository to get estate stats
	estate, err := s.Repository.GetDetailEstate(ctx.Request().Context(), repository.GetDetailEstateInput{
		ID: id,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Failed to get estate details",
		})
	}
	fmt.Println("Estate:", estate.Trees)
	// Calculate stats
	count := len(estate.Trees)
	min := int(estate.Trees[0].Height)
	max := int(estate.Trees[0].Height)
	median := int(estate.Trees[0].Height)
	for _, tree := range estate.Trees {
		if int(tree.Height) < min {
			min = int(tree.Height)
		}
		if int(tree.Height) > max {
			max = int(tree.Height)
		}
		median += int(tree.Height)
	}
	median /= count
	return ctx.JSON(http.StatusOK, generated.EstateStatsResponse{
		"count":  count,
		"min":    min,
		"max":    max,
		"median": median,
	})
}

func (s *Server) GetEstateIdDronePlan(ctx echo.Context, id types.UUID) error {
	// get detail estate
	estate, err := s.Repository.GetDetailEstate(ctx.Request().Context(), repository.GetDetailEstateInput{
		ID: id,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Failed to get estate details",
		})
	}

	length := int(estate.Length)
	width := int(estate.Width)
	// crerate a slice/list of trees
	trees := []Tree{}
	for _, tree := range estate.Trees {
		trees = append(trees, Tree{
			X:      int(tree.X),
			Y:      int(tree.Y),
			Height: int(tree.Height),
		})
	}

	totalDistance := CalculateDroneDistance(length, width, trees)
	return ctx.JSON(http.StatusOK, generated.TotalDistanceResponse{
		TotalDistance: &totalDistance,
	})
}

func (s *Server) GetEstateIdDronePlanMax(ctx echo.Context, id types.UUID, params generated.GetEstateIdDronePlanMaxParams) error {
	max := params.MaxDistance

	estate, err := s.Repository.GetDetailEstate(ctx.Request().Context(), repository.GetDetailEstateInput{
		ID: id,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Failed to get estate details",
		})
	}

	length := int(estate.Length)
	width := int(estate.Width)
	// crerate a slice/list of trees
	trees := []Tree{}
	for _, tree := range estate.Trees {
		trees = append(trees, Tree{
			X:      int(tree.X),
			Y:      int(tree.Y),
			Height: int(tree.Height),
		})
	}

	lastRest := MaxDistanceDrone(length, width, trees, max)
	if lastRest == nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid max distance",
		})
	}
	if len(lastRest) != 2 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid last rest coordinates",
		})
	}
	// Check if the last rest coordinates are within the estate bounds
	if lastRest[0] < 1 || lastRest[0] > length || lastRest[1] < 1 || lastRest[1] > width {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: fmt.Sprintf("Last rest coordinates (%d, %d) are out of bounds for estate (%d, %d)", lastRest[0], lastRest[1], width, length),
		})
	}
	return ctx.JSON(http.StatusOK, generated.DronePlanResponse{
		MaxDistance: max,
		Rest: generated.DroneRest{
			X: &lastRest[0],
			Y: &lastRest[1],
		},
	})
}
