package repository

import (
	"context"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) CreateEstate(ctx context.Context, input CreateEstateInput) (output CreateEstateOutput, err error) {
	// Insert the new estate
	err = r.Db.QueryRowContext(ctx, "INSERT INTO estate ( width, length) VALUES ( $1, $2) RETURNING id_estate", input.Width, input.Length).Scan(&output.ID)
	if err != nil {
		return
	}
	return
}

func (r *Repository) CreateTree(ctx context.Context, input CreateTreeInput) (output CreateTreeOutput, err error) {
	// Insert the new tree
	err = r.Db.QueryRowContext(ctx, "INSERT INTO tree ( id_estate, x, y) VALUES ($1, $2, $3) RETURNING id_tree", input.EstateID, input.X, input.Y).Scan(&output.ID)
	if err != nil {
		return
	}
	return
}
