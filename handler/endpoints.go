package handler

import (
	"fmt"
	"net/http"

	"github.com/candrap89/droneTravelingDistance/generated"
	"github.com/candrap89/droneTravelingDistance/repository"
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
