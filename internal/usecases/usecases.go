package usecases

import (
	"context"
	"youGoToPairs/internal/models"
)

type ScheduleUseCase interface {
	AddSchedule(ctx context.Context, groupId int64, schedule models.Schedule) error
	GetSchedule(ctx context.Context, groupId int64) (models.Schedule, error)
	StartSendPolls()
}
