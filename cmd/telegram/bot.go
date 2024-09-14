package main

import (
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
