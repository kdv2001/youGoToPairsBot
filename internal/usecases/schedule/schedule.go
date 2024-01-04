package schedule

import (
	"context"
	"errors"
	"fmt"
	"time"
	"youGoToPairs/internal/models"
	"youGoToPairs/internal/repo"
)

var days = map[string]time.Weekday{
	time.Monday.String():    time.Monday,
	time.Thursday.String():  time.Thursday,
	time.Saturday.String():  time.Saturday,
	time.Wednesday.String(): time.Wednesday,
	time.Friday.String():    time.Friday,
	time.Tuesday.String():   time.Tuesday,
	time.Sunday.String():    time.Sunday,
}

const secondInDay = 86400

type UseCase struct {
	scheduleRepo repo.ScheduleRepo
}

func NewScheduleUseCase(r repo.ScheduleRepo) *UseCase {
	return &UseCase{
		scheduleRepo: r,
	}
}

func (u *UseCase) AddSchedule(ctx context.Context, groupId int64, schedule models.Schedule) error {
	keys := make([]int64, 0, len(schedule.Days))
	addSchedules := make(map[string]models.DaySchedule)
	for _, val := range schedule.Days {
		multipluer, ok := days[val.DayOfTheWeek]
		if !ok {
			return errors.New("bad day")

		}

		keys = append(keys, int64(secondInDay*int(multipluer)+int(val.Time)))
		key := fmt.Sprint(secondInDay*int(multipluer) + int(val.Time))

		addSchedules[key] = val
	}

	if err := u.scheduleRepo.AddSchedule(ctx, groupId, addSchedules); err != nil {
		return err
	}

	if err := u.scheduleRepo.BulkSetShedule(ctx, groupId, keys); err != nil {
		return err
	}

	return nil
}

func (u *UseCase) GetSchedule(ctx context.Context, groupId int64) (models.Schedule, error) {
	res, err := u.scheduleRepo.GetAllSchedules(ctx, 0, 10000000)
	if err != nil {
		return models.Schedule{}, err
	}

	fmt.Println(res)

	// days := make([]models.DaySchedule, 0, len(res))
	// for _, val := range res {
	// 	days = append(days, val)
	// }

	return models.Schedule{}, nil
}

func (u *UseCase) StartSendPolls() {
	go u.sendPolls(context.Background())
}

func (u *UseCase) sendPolls(ctx context.Context) {
	period := 60 * time.Second

	ticker := time.NewTicker(period)
	for _ = range ticker.C {
		timeNow := time.Now()
		curentDaySeconds := time.Date(0, 0, 0, timeNow.Hour(), timeNow.Minute(), timeNow.Second(), 0, time.UTC).Unix()

		minTime := secondInDay*int64(timeNow.Weekday()) + curentDaySeconds
		maxTime := secondInDay*int64(timeNow.Weekday()) + curentDaySeconds + int64(period.Seconds())

		schedules, err := u.scheduleRepo.GetAllSchedules(ctx, minTime, maxTime)
		if err != nil {
			continue
		}

		fmt.Println(schedules)

		for groupId, schedule := range schedules {
			for _, curTime := range schedule {
				u.scheduleRepo.GetSchedule(ctx, groupId)

			}
		}
	}
}
