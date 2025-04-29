// services/report_service.go
package services

import (
	"context"
	"errors"
	"report-service/internal/entities"
	"report-service/internal/repositories"
)

type ReportService interface {
	ProcessReport(ctx context.Context, report *entities.Report) (string, error)
}

type reportService struct {
	repo repositories.ReportRepository
}

func NewReportService(repo repositories.ReportRepository) *reportService {
	return &reportService{repo: repo}
}

func (s *reportService) ProcessReport(ctx context.Context, report *entities.Report) (string, error) {
	// Validar reporte
	if report.WebID == "" || len(report.Elements) == 0 {
		return "", errors.New("datos del reporte inválidos")
	}

	// Guardar reporte
	if err := s.repo.SaveReport(ctx, report); err != nil {
		return "", err
	}

	// Obtener conteo
	count, err := s.repo.GetReportCount(ctx, report.WebID)
	if err != nil {
		return "", err
	}

	// Verificar si es el tercer reporte
	if count >= 3 {
		return "Se ha reportado el número máximo de materiales, nuestro equipo se está encargando", nil
	}

	return "Reporte recibido correctamente", nil
}
