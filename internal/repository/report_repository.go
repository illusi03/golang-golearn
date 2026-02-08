package repository

import (
	"context"
	"time"

	"github.com/illusi03/golearn/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepository struct {
	dbPool *pgxpool.Pool
}

func NewReportRepository(dbPool *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{dbPool: dbPool}
}

func (r *ReportRepository) GetReport(ctx context.Context, startDate, endDate time.Time) (*model.ReportModel, error) {
	report := &model.ReportModel{}

	const summaryQuery = `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COUNT(*) as total_transaction
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`
	err := r.dbPool.QueryRow(ctx, summaryQuery, startDate, endDate).Scan(
		&report.TotalRevenue,
		&report.TotalTransaction,
	)
	if err != nil {
		return nil, err
	}

	const bestSellerQuery = `
		SELECT 
			p.name,
			COALESCE(SUM(td.quantity), 0) as qty_sold
		FROM transaction_details td
		JOIN transactions t ON t.id = td.transaction_id
		JOIN products p ON p.id = td.product_id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.id, p.name
		ORDER BY qty_sold DESC
		LIMIT 1
	`
	var bestSeller model.BestSellerModel
	err = r.dbPool.QueryRow(ctx, bestSellerQuery, startDate, endDate).Scan(
		&bestSeller.Name,
		&bestSeller.QtySold,
	)
	if err == nil {
		report.BestSeller = &bestSeller
	}

	return report, nil
}
