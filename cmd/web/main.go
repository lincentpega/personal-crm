package main

import (
	"context"
	"html/template"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/lincentpega/personal-crm/internal/config"
	"github.com/lincentpega/personal-crm/internal/db"
	"github.com/lincentpega/personal-crm/internal/log"
)

type application struct {
	log            *log.Logger
	templates      map[string]*template.Template
	sessionManager *scs.SessionManager
}

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

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(database)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		log:            log,
		sessionManager: sessionManager,
	}

	if err := app.loadTemplates(); err != nil {
		log.ErrorLog.Fatal(err)
	}

	s := &http.Server{
		Addr:     config.Addr,
		Handler:  app.route(),
		ErrorLog: log.ErrorLog,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.InfoLog.Printf("Starting server on %s", config.Addr)
		err := s.ListenAndServe()
		log.ErrorLog.Fatal(err)
	}()

	<-ctx.Done()
	log.InfoLog.Println("Shutting down server")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.Shutdown(ctx)
	log.ErrorLog.Fatal(err)
}
