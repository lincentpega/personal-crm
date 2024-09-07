package main

import (
	tele "gopkg.in/telebot.v3"
)

func (app *Application) handle(bot *tele.Bot)  {
    bot.Handle("/hello", func(ctx tele.Context) error {
        return ctx.Send("Hello, world!")
    })
}
