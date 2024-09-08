package main

import (
	_ "github.com/lib/pq"
	"github.com/lincentpega/personal-crm/internal/config"
	"github.com/lincentpega/personal-crm/internal/db"
	"github.com/lincentpega/personal-crm/internal/log"
	"github.com/lincentpega/personal-crm/internal/models"
)


func main() {
    config := config.Load()
	log := log.New()

    db, err := db.Connect(config.DSN)
    if err != nil {
        log.ErrorLog.Fatal(err)
    }

	_ = models.NewPersonModel(db)

    b, err := newBot(config.Token, log)
    if err != nil {
        log.ErrorLog.Fatal(err)
    }

    name, err := b.name()
    if err != nil {
        name = ""
    }

	log.InfoLog.Printf("Starting bot %s", name)
    b.start()
}
