package main

import (
	"context"
	"time"

	"github.com/lincentpega/personal-crm/internal/log"
	"github.com/lincentpega/personal-crm/internal/models/notifications"
	"github.com/lincentpega/personal-crm/internal/models/person"
	"gopkg.in/telebot.v3"
)

type bot struct {
	*telebot.Bot
	personRepo *person.PersonRepository
	notifRepo  *notifications.NotificationRepository
	log        *log.Logger
}

func newBot(token string, log *log.Logger, pr *person.PersonRepository, nr *notifications.NotificationRepository) (*bot, error) {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &bot{Bot: b, log: log, personRepo: pr, notifRepo: nr}, nil
}

func (b *bot) logStart() error {

	botInfo, err := b.MyName("")
	if err != nil {
		return err
	}

	b.log.InfoLog.Printf("Starting bot %s", botInfo.Name)

	return nil
}

func (b *bot) route() {
	base := b.Group()

	base.Handle("/create-person", func(ctx telebot.Context) error {
		firstName := "Igor"
		lastName := "Krasnyukov"
		err := b.personRepo.Insert(
			context.Background(),
			&person.Person{
				FirstName: firstName,
				LastName:  &lastName,
			})
		if err != nil {
			return err
		}
		return ctx.Send("Person is created")
	})

	base.Handle("/create-notification", func(ctx telebot.Context) error {
		now := time.Now().UTC()
		err := b.notifRepo.Insert(
			context.Background(),
			&notifications.Notification{
				PersonID:         2,
				NotificationTime: &now,
				Status:           notifications.Pending,
				Type:             notifications.KeepInTouch,
			})
		if err != nil {
			return err
		}
		return ctx.Send("Notification scheduled")
	})

	base.Handle("/hello", func(ctx telebot.Context) error {
		var kbd [][]telebot.InlineButton
		btn1 := telebot.InlineButton{Text: "SOSAT", Data: "sosat"}
		row1 := []telebot.InlineButton{btn1}
		kbd = append(kbd, row1)
		mrkp := &telebot.ReplyMarkup{InlineKeyboard: kbd}
		return ctx.Send("Hello, world!", mrkp)
	})

	base.Handle(telebot.OnCallback, func(ctx telebot.Context) error {
		c := ctx.Callback()
		return ctx.Send(c.Data)
	})
}
