package telegramtelebot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"youGoToPairs/internal/models"
	"youGoToPairs/internal/usecases"

	"gopkg.in/telebot.v3"
)

type Handlers struct {
	scheduleUC usecases.ScheduleUseCase
}

func NewScheduleHandlers(s usecases.ScheduleUseCase) *Handlers {
	return &Handlers{
		scheduleUC: s,
	}
}

func (h *Handlers) AddSchedule(ctx telebot.Context) error {
	fmt.Println(ctx.Message().Text)

	res := strings.Split(ctx.Message().Text, "@")
	fmt.Println(res)

	if len(res) != 2 {
		return errors.New("bad config")
	}

	scheduleModel := models.Schedule{}
	err := json.Unmarshal([]byte(res[1]), &scheduleModel)
	if err != nil {
		return err
	}

	err = h.scheduleUC.AddSchedule(context.Background(), ctx.Chat().ID, scheduleModel)
	if err != nil {
		log.Fatal(err)
	}

	return ctx.Send(fmt.Sprint(ctx.Chat().ID))
}

func (h *Handlers) GetSchedule(ctx telebot.Context) error {

	h.scheduleUC.GetSchedule(context.Background(), 123)

	return ctx.Send("fdf", &telebot.SendOptions{})
}
