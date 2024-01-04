package main

import (
	"context"
	"fmt"
	"log"

	scheduleRedisRepo "youGoToPairs/internal/repo/redis"

	"github.com/redis/go-redis/v9"
)

func main() {
	// bot, err := app.CreateApp()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bot.Run()

	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	scheduleRepo := scheduleRedisRepo.NewScheduleRepo(rdb)
	scheduleRepo.BulkSetShedule(context.Background(), 123, []int64{123, 1234, 1235})
	res, err := scheduleRepo.GetAllSchedules(context.Background(), 0, 10000000000)
	fmt.Println(res, err)
}
