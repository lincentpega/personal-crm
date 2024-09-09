package main

import (
	"time"

	"github.com/lincentpega/personal-crm/internal/log"
	tele "gopkg.in/telebot.v3"
)

type bot struct {
	*tele.Bot
	log *log.Logger
}

func newBot(token string, log *log.Logger) (*bot, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
        return nil, err
	}

    return &bot{Bot: b, log: log}, nil
}

func (b *bot) name() (string, error) {
	botInfo, err := b.MyName("")
	if err != nil {
        return "", err
	}

    return botInfo.Name, nil
}

func (b *bot) route() {
    base := b.Group()

    base.Handle("/hello", func(ctx tele.Context) error {
        return ctx.Send("Hello, world!")
    })
}

func (b *bot) start() {
    b.route()
    b.Start()
}

