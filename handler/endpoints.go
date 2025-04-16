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
	return ctx.JSON(http.StatusCreated, generated.CreateEstateResponse{
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
	if input.X < 0 || input.Y < 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "X and Y coordinates must be greater than or equal to 0",
		})
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
	return ctx.JSON(http.StatusCreated, generated.CreateTreeResponse{
		Id: &output.ID,
	})
}
