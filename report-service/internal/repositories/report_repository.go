// repositories/report_repository.go
package repositories

import (
	"context"
	"report-service/internal/entities"

	"github.com/go-redis/redis/v8"
)

type ReportRepository interface {
	SaveReport(ctx context.Context, report *entities.Report) error
	GetReportCount(ctx context.Context, webID string) (int, error)
}

type RedisReportRepository struct {
	client *redis.Client
}

func NewRedisReportRepository(client *redis.Client) *RedisReportRepository {
	return &RedisReportRepository{client: client}
}

func (r *RedisReportRepository) SaveReport(ctx context.Context, report *entities.Report) error {
	// Incrementar contador para este WebID
	countKey := "report_count:" + report.WebID
	_, err := r.client.Incr(ctx, countKey).Result()
	return err
}

func (r *RedisReportRepository) GetReportCount(ctx context.Context, webID string) (int, error) {
	countKey := "report_count:" + webID
	count, err := r.client.Get(ctx, countKey).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}
