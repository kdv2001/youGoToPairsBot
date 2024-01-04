package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"youGoToPairs/internal/models"

	"github.com/redis/go-redis/v9"
)

const (
	groupScheduleKey = "groupSchedule"
	scheduleKey      = "schedule"
)

type Repo struct {
	r *redis.Client
}

func NewScheduleRepo(client *redis.Client) *Repo {
	return &Repo{
		r: client,
	}
}

func (r *Repo) AddSchedule(ctx context.Context, groupId int64, schedule map[string]models.DaySchedule) error {
	data := make([]interface{}, 0, len(schedule))
	for timeKey, val := range schedule {
		byteVal, err := json.Marshal(val)
		if err != nil {
			return err
		}

		data = append(data, timeKey, byteVal)
	}

	err := r.r.HSet(ctx, createScheduleKey(groupId), data...).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetSchedule(ctx context.Context, groupId int64, timeKey int) (models.DaySchedule, error) {
	result, err := r.r.HGet(ctx, createScheduleKey(groupId), fmt.Sprint(timeKey)).Result()
	if err != nil {
		return models.DaySchedule{}, err
	}

	scheduleVal := models.DaySchedule{}
	err = json.Unmarshal([]byte(result), &scheduleVal)
	if err != nil {
		return scheduleVal, errors.New("bad get result")
	}

	return scheduleVal, nil
}

func (r *Repo) GetSchedules(ctx context.Context, groupId int64) (map[string]models.DaySchedule, error) {
	getResult, err := r.r.HGetAll(ctx, createScheduleKey(groupId)).Result()
	if err != nil {
		return nil, err
	}

	scheduleModel := make(map[string]models.DaySchedule)
	for key, val := range getResult {
		scheduleVal := models.DaySchedule{}
		err = json.Unmarshal([]byte(val), &scheduleVal)
		if err != nil {
			return nil, errors.New("bad get result")
		}

		scheduleModel[key] = scheduleVal
	}

	return scheduleModel, nil
}

func (r *Repo) BulkSetShedule(ctx context.Context, groupId int64, keys []int64) error {
	members := make([]redis.Z, 0, len(keys))
	for _, key := range keys {
		members = append(members, redis.Z{
			Score:  float64(key),
			Member: fmt.Sprintf("%d:%d", groupId, key),
		})

	}

	err := r.r.ZAdd(ctx, fmt.Sprint(scheduleKey, ":", groupId), members...).Err()
	return err
}

func (r *Repo) BulkDelShedule(ctx context.Context, groupId int64, keys []int64) error {
	pipe := r.r.Pipeline()

	for _, key := range keys {
		pipe.Del(ctx, fmt.Sprint(scheduleKey, ":", groupId, ":", key))
	}

	_, err := pipe.Exec(ctx)

	return err
}

func (r *Repo) GetAllSchedules(ctx context.Context, minTime, maxTime int64) (map[int64][]int64, error) {
	res, err := r.r.Keys(ctx, fmt.Sprint("*", scheduleKey, "*")).Result()
	if err != nil {
		return nil, err
	}

	pipe := r.r.Pipeline()

	scheduleRes := make(map[int64][]int64)
	for _, key := range res {
		splitKey := strings.Split(key, ":")
		if len(splitKey) != 2 {
			return nil, errors.New("bad key")
		}

		groupId, err := strconv.ParseInt(splitKey[1], 10, 64)
		if err != nil {
			return nil, err
		}

		pipe.ZRangeByScore(ctx, fmt.Sprint(scheduleKey, ":", groupId), &redis.ZRangeBy{
			Min: strconv.Itoa(int(minTime)),
			Max: strconv.Itoa(int(maxTime)),
		})
	}

	execRes, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	for _, val := range execRes {
		getMember := val.(*redis.StringSliceCmd)
		fmt.Println(getMember)
		for _, curMember := range getMember.Val() {
			splitKey := strings.Split(curMember, ":")
			if len(splitKey) != 2 {
				return nil, errors.New("bad key")
			}

			groupId, err := strconv.ParseInt(splitKey[0], 10, 64)
			if err != nil {
				return nil, err
			}

			intCurScore, err := strconv.ParseInt(splitKey[1], 10, 64)
			if err != nil {
				return nil, err
			}

			scheduleRes[groupId] = append(scheduleRes[groupId], intCurScore)
		}
	}

	return scheduleRes, nil
}

func createScheduleKey(groupId int64) string {
	return fmt.Sprint(groupScheduleKey, ":", groupId)
}
