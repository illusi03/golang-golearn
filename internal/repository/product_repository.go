package repository

import (
	"context"
	"errors"

	"github.com/illusi03/golearn/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	dbPool *pgxpool.Pool
}

func NewProductRepository(dbPool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		dbPool: dbPool,
	}
}

func (a *ProductRepository) FindAll(ctx context.Context, name string) ([]model.ProductModel, error) {
	query := `
		SELECT p.id, p.name, p.description, p.price, p.stock, p.category_id, c.name as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
	`
	args := []any{}

	if name != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	query += " ORDER BY p.id"

	rows, err := a.dbPool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.ProductModel, 0)
	for rows.Next() {
		var c model.ProductModel
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Price, &c.Stock, &c.CategoryID, &c.CategoryName); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (a *ProductRepository) FindOne(
	ctx context.Context,
	id int,
) (*model.ProductModel, error) {
	const query = `
		SELECT p.id, p.name, p.description, p.price, p.stock, p.category_id, c.name as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`
	var c model.ProductModel
	err := a.dbPool.
		QueryRow(ctx, query, id).
		Scan(&c.ID, &c.Name, &c.Description, &c.Price, &c.Stock, &c.CategoryID, &c.CategoryName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (a *ProductRepository) Delete(
	ctx context.Context,
	id int,
) (bool, error) {
	const query = `
		DELETE FROM products
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

func (a *ProductRepository) Update(
	ctx context.Context,
	c *model.ProductModel,
) (bool, error) {
	const query = `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock = $4, category_id = $5
		WHERE id = $6
	`
	cmdTag, err := a.dbPool.Exec(
		ctx,
		query,
		c.Name,
		c.Description,
		c.Price,
		c.Stock,
		c.CategoryID,
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

func (a *ProductRepository) Create(
	ctx context.Context,
	c *model.ProductModel,
) (*model.ProductModel, error) {
	const query = `
		INSERT INTO products (name, description, price, stock, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, price, stock, category_id
	`
	var out model.ProductModel
	err := a.dbPool.QueryRow(
		ctx,
		query,
		c.Name,
		c.Description,
		c.Price,
		c.Stock,
		c.CategoryID,
	).Scan(&out.ID, &out.Name, &out.Description, &out.Price, &out.Stock, &out.CategoryID)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
