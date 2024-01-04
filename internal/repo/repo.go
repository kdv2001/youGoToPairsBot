package repo

import (
	"context"
	"youGoToPairs/internal/models"
)

type ScheduleRepo interface {
	AddSchedule(ctx context.Context, groupId int64, schedule map[string]models.DaySchedule) error
	GetSchedule(ctx context.Context, groupId int64, getTime int)
	GetSchedules(ctx context.Context, groupId int64) (map[string]models.DaySchedule, error)
	BulkSetShedule(ctx context.Context, groupId int64, keys []int64) error
	BulkDelShedule(ctx context.Context, groupId int64, keys []int64) error
	GetAllSchedules(ctx context.Context, maxTime, minTime int64) (map[int64][]int64, error)
}
