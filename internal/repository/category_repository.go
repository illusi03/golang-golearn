package repository

import (
	"context"
	"errors"

	"github.com/illusi03/golearn/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	dbPool *pgxpool.Pool
}

func NewCategoryRepository(dbPool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		dbPool: dbPool,
	}
}

func (a *CategoryRepository) FindAll(ctx context.Context) ([]model.CategoryModel, error) {
	const query = `
		SELECT id, name, description
		FROM categories
		ORDER BY id
	`
	rows, err := a.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.CategoryModel, 0)
	for rows.Next() {
		var c model.CategoryModel
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (a *CategoryRepository) FindOne(
	ctx context.Context,
	id int,
) (*model.CategoryModel, error) {
	const query = `
		SELECT id, name, description
		FROM categories
		WHERE id = $1
	`
	var c model.CategoryModel
	err := a.dbPool.
		QueryRow(ctx, query, id).
		Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (a *CategoryRepository) Delete(
	ctx context.Context,
	id int,
) (bool, error) {
	const query = `
		DELETE FROM categories
		WHERE id = $1
	`
	cmdTag, err := a.dbPool.Exec(ctx, query, id)
	if err != nil {
		return false, err
	}
	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}
	return true, nil
}

func (a *CategoryRepository) Update(
	ctx context.Context,
	c *model.CategoryModel,
) (bool, error) {
	const query = `
		UPDATE categories
		SET name = $1, description = $2
		WHERE id = $3
	`
	cmdTag, err := a.dbPool.Exec(
		ctx,
		query,
		c.Name,
		c.Description,
		c.ID,
	)
	if err != nil {
		return false, err
	}
	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}
	return true, nil
}

func (a *CategoryRepository) Create(
	ctx context.Context,
	c *model.CategoryModel,
) (*model.CategoryModel, error) {
	const query = `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description
	`
	var out model.CategoryModel
	err := a.dbPool.QueryRow(
		ctx,
		query,
		c.Name,
		c.Description,
	).Scan(&out.ID, &out.Name, &out.Description)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
