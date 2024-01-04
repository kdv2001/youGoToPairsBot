package app

import (
	"context"
	"os"
	"time"

	telegramtelebot "youGoToPairs/internal/handlers/telegramTelebot"
	scheduleRedisRepo "youGoToPairs/internal/repo/redis"
	scheduleUseCase "youGoToPairs/internal/usecases/schedule"

	"github.com/redis/go-redis/v9"
	"gopkg.in/telebot.v3"
)

type app struct {
	bot *telebot.Bot
}

func CreateApp() (*app, error) {
	settings := telebot.Settings{
		Token:  os.Getenv("botKey"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	scheduleRepo := scheduleRedisRepo.NewScheduleRepo(rdb)
	scheduleUC := scheduleUseCase.NewScheduleUseCase((scheduleRepo))
	scheduleHandlers := telegramtelebot.NewScheduleHandlers(scheduleUC)

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	bot.Handle("/addSchedule", scheduleHandlers.AddSchedule)
	bot.Handle("/getSchedule", scheduleHandlers.GetSchedule)

	return &app{
		bot: bot,
	}, nil
}

func (a app) Run() {
	a.bot.Start()
}
