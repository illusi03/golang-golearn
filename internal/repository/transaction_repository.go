package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/illusi03/golearn/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	dbPool *pgxpool.Pool
}

func NewTransactionRepository(dbPool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		dbPool: dbPool,
	}
}

func (r *TransactionRepository) Create(
	ctx context.Context,
	transaction *model.TransactionModel,
) (*model.TransactionModel, error) {
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Insert transaction
	const txQuery = `
		INSERT INTO transactions (total_amount)
		VALUES ($1)
		RETURNING id, total_amount, created_at
	`
	err = tx.QueryRow(ctx, txQuery, transaction.TotalAmount).Scan(
		&transaction.ID,
		&transaction.TotalAmount,
		&transaction.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Batch insert transaction details
	if len(transaction.Details) > 0 {
		valueStrings := make([]string, len(transaction.Details))
		args := make([]interface{}, 0, len(transaction.Details)*4)

		for i := range transaction.Details {
			transaction.Details[i].TransactionID = transaction.ID
			detail := &transaction.Details[i]
			offset := i * 4
			valueStrings[i] = fmt.Sprintf("($%d, $%d, $%d, $%d)", offset+1, offset+2, offset+3, offset+4)
			args = append(args, detail.TransactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
		}

		detailQuery := fmt.Sprintf(`
			INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal)
			VALUES %s
			RETURNING id
		`, strings.Join(valueStrings, ", "))

		rows, err := tx.Query(ctx, detailQuery, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		i := 0
		for rows.Next() {
			if err := rows.Scan(&transaction.Details[i].ID); err != nil {
				return nil, err
			}
			i++
		}

		// Batch deduct stock from products
		stockUpdateStrings := make([]string, len(transaction.Details))
		stockArgs := make([]interface{}, 0, len(transaction.Details)*2)

		for i, detail := range transaction.Details {
			offset := i * 2
			stockUpdateStrings[i] = fmt.Sprintf("($%d::int, $%d::int)", offset+1, offset+2)
			stockArgs = append(stockArgs, detail.ProductID, detail.Quantity)
		}

		stockQuery := fmt.Sprintf(`
			UPDATE products AS p
			SET stock = p.stock - v.quantity
			FROM (VALUES %s) AS v(product_id, quantity)
			WHERE p.id = v.product_id AND p.stock >= v.quantity
		`, strings.Join(stockUpdateStrings, ", "))

		result, err := tx.Exec(ctx, stockQuery, stockArgs...)
		if err != nil {
			return nil, err
		}

		if result.RowsAffected() != int64(len(transaction.Details)) {
			return nil, fmt.Errorf("insufficient stock: concurrent modification detected")
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *TransactionRepository) FindAll(
	ctx context.Context,
) ([]model.TransactionModel, error) {
	const query = `
		SELECT 
			t.id, t.total_amount, t.created_at,
			td.id, td.product_id, p.name, td.quantity, td.subtotal
		FROM transactions t
		LEFT JOIN transaction_details td ON td.transaction_id = t.id
		LEFT JOIN products p ON p.id = td.product_id
		ORDER BY t.created_at DESC, td.id ASC
	`
	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txMap := make(map[int]*model.TransactionModel)
	var txOrder []int

	for rows.Next() {
		var txID int
		var totalAmount int
		var createdAt time.Time
		var detailID *int
		var productID *int
		var productName *string
		var quantity *int
		var subtotal *int

		err = rows.Scan(
			&txID, &totalAmount, &createdAt,
			&detailID, &productID, &productName, &quantity, &subtotal,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := txMap[txID]; !exists {
			txMap[txID] = &model.TransactionModel{
				ID:          txID,
				TotalAmount: totalAmount,
				CreatedAt:   createdAt,
				Details:     []model.TransactionDetailModel{},
			}
			txOrder = append(txOrder, txID)
		}

		if detailID != nil {
			txMap[txID].Details = append(txMap[txID].Details, model.TransactionDetailModel{
				ID:            *detailID,
				TransactionID: txID,
				ProductID:     *productID,
				ProductName:   productName,
				Quantity:      *quantity,
				Subtotal:      *subtotal,
			})
		}
	}

	transactions := make([]model.TransactionModel, 0, len(txOrder))
	for _, id := range txOrder {
		transactions = append(transactions, *txMap[id])
	}

	return transactions, nil
}
