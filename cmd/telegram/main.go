package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v3"
)

type Application struct {
    DB *sql.DB
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func (app *Application) botName(b *tele.Bot) string {
	botInfo, err := b.MyName("")
	if err != nil {
		app.ErrorLog.Fatal(err)
	}

	return botInfo.Name
}

func (app *Application) dbConnect(dsn string) *sql.DB {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        app.ErrorLog.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        app.ErrorLog.Fatal(err)
    }

    return db
}

func main() {
	token := flag.String("token", "empty_token", "telegram bot token")
    dsn := flag.String("dsn", "host=localhost port=5433 user=postgres password=mysecretpassword dbname=postgres sslmode=disable", "PostgreSQL datasource name")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

    app.dbConnect(*dsn)

	pref := tele.Settings{
		Token:  *token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		errorLog.Fatal(err)
		return
	}

	app.handle(b)

    botName := app.botName(b)

	infoLog.Printf("Starting bot %s", botName)
	b.Start()
}
