package main

import (
	"context"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/lincentpega/personal-crm/internal/config"
	"github.com/lincentpega/personal-crm/internal/db"
	"github.com/lincentpega/personal-crm/internal/log"
)


func main() {
    config := config.Load()
	log := log.New()

    database, err := db.Connect(config.DSN)
    if err != nil {
        log.ErrorLog.Fatal(err)
    }
    defer database.Close()

    err = db.ExecMigrations(database, log)
    if err != nil {
        log.ErrorLog.Fatal(err)
    }

    b, err := newBot(config.Token, log)
    if err != nil {
        log.ErrorLog.Fatal(err)
    }

    name, err := b.name()
    if err != nil {
        name = ""
    }

	log.InfoLog.Printf("Starting bot %s", name)
    go b.start()

    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()
    <-ctx.Done()

    log.InfoLog.Println("Shutting down bot")
    b.Stop()
}
