// This file contains types that are used in the repository layer.
package repository

import uuid "github.com/google/uuid"

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type CreateEstateInput struct {
	Width  int // Width of the estate
	Length int // Length of the estate
}
type CreateEstateOutput struct {
	ID uuid.UUID
}

type CreateTreeInput struct {
	EstateID uuid.UUID
	X        int // X coordinate of the tree
	Y        int // Y coordinate of the tree
	Height   int // Height of the tree
}
type CreateTreeOutput struct {
	ID uuid.UUID
}
