package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// DoWithRetry makes an HTTP request with retry logic.
// method: "GET", "POST", etc.
// url: the request URL
// retries: how many times to retry
// backoff: delay between retries (grows if exponential=true)
// exponential: whether to use exponential backoff
func DoWithRetry(ctx context.Context, method string, url string, retries int, backoff time.Duration, exponential bool) (*http.Response, error) {
	var resp *http.Response
	var err error

	client := http.DefaultClient

	for i := 0; i < retries; i++ {
		req, reqErr := http.NewRequestWithContext(ctx, method, url, nil)
		if reqErr != nil {
			return nil, fmt.Errorf("failed to create request: %w", reqErr)
		}

		resp, err = client.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Success
			return resp, nil
		}

		if resp != nil {
			resp.Body.Close() // prevent resource leak
		}

		fmt.Printf("Attempt %d failed: %v\n", i+1, err)
		time.Sleep(backoff)

		if exponential {
			backoff *= 2
		}
	}

	return nil, errors.New("request failed after retries")
}

// Example usage
// func main() {
// 	ctx := context.Background()
// 	url := "https://jsonmock.hackerrank.com/api/food_outlets?city=New%20York&estimated_cost=100"
// 	resp, err := DoWithRetry(ctx, "GET", url, 3, 2*time.Second, true)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	defer resp.Body.Close()
//
