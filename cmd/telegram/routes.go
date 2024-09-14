package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/lincentpega/personal-crm/internal/models/notifications"
	"github.com/lincentpega/personal-crm/internal/models/person"
	"gopkg.in/telebot.v3"
)

func (b *bot) route() {
	base := b.Group()

	base.Handle("/create-person", func(ctx telebot.Context) error {
		firstName := "Igor"
		lastName := "Krasnyukov"
		err := b.personRepo.Insert(
			context.Background(),
			&person.Person{
				FirstName: firstName,
				LastName:  sql.NullString{String: lastName, Valid: true},
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
				NotificationTime: now,
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
