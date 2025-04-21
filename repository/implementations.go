package repository

import (
	"context"
	"fmt"
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
	err = r.Db.QueryRowContext(ctx, "INSERT INTO tree ( id_estate, x, y, height) VALUES ($1, $2, $3, $4) RETURNING id_tree", input.EstateID, input.X, input.Y, input.Height).Scan(&output.ID)
	if err != nil {
		return
	}
	return
}
func (r *Repository) GetDetailEstate(ctx context.Context, input GetDetailEstateInput) (output GetDetailEstateOutput, err error) {
	// Get the estate details
	// err = r.Db.QueryRowContext(ctx, "SELECT id_estate, width, length FROM estate WHERE id_estate = $1", input.ID).Scan(&output.ID, &output.Width, &output.Length)
	// if err != nil {
	// 	return
	// }
	// Get the trees associated with the estate
	rows, err := r.Db.QueryContext(ctx, "SELECT id_tree, x, y, height FROM tree WHERE id_estate = $1", input.ID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tree Tree
		err := rows.Scan(&tree.ID, &tree.X, &tree.Y, &tree.Height)
		// Check for errors from the row scan
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return output, err
		}
		output.Trees = append(output.Trees, tree)
	}
	return
}
