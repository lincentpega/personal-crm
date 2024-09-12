package main

import (
	"context"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/lincentpega/personal-crm/internal/config"
	"github.com/lincentpega/personal-crm/internal/db"
	"github.com/lincentpega/personal-crm/internal/log"
	"github.com/lincentpega/personal-crm/internal/models/notifications"
	"github.com/lincentpega/personal-crm/internal/models/person"
	"github.com/lincentpega/personal-crm/internal/services"
)

func main() {
	config := config.Load()
	log := log.New()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	database, err := db.Connect(config.DSN)
	if err != nil {
		log.ErrorLog.Fatal(err)
	}
	defer database.Close()

	err = db.ExecMigrations(database, log)
	if err != nil {
		log.ErrorLog.Fatal(err)
	}

	notificaitonRepo := notifications.NewRepository(database)
	personRepo := person.NewRepository(database)

	b, err := newBot(config.Token, log, personRepo, notificaitonRepo)
	if err != nil {
		log.ErrorLog.Fatal(err)
	}

	notificationService := services.NewNotificationService(b.Bot, notificaitonRepo, personRepo, log, config)

	startApplication(ctx, b, notificationService)

	<-ctx.Done()

	log.InfoLog.Println("Shutting down bot")
	b.Stop()
}

func startApplication(ctx context.Context, b *bot, ns *services.NotificationService) {
	b.route()
	b.logStart()
	go b.Start()
	go ns.ProcessNotifications(ctx)
}
