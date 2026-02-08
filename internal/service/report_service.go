package service

import (
	"context"
	"time"

	"github.com/illusi03/golearn/internal/model"
	"github.com/illusi03/golearn/internal/repository"
)

type ReportService struct {
	reportRepository *repository.ReportRepository
}

func NewReportService(reportRepository *repository.ReportRepository) *ReportService {
	return &ReportService{reportRepository: reportRepository}
}

func (s *ReportService) GetReport(ctx context.Context, startDate, endDate time.Time) (*model.ReportModel, error) {
	return s.reportRepository.GetReport(ctx, startDate, endDate)
}
