package config

import "flag"

type AppConfig struct {
	Token  string
	DSN    string
	UserID int
}

func Load() *AppConfig {
	token := flag.String("token", "empty_token", "telegram bot token")
	dsn := flag.String("dsn", "host=localhost port=5433 user=postgres password=mysecretpassword dbname=postgres sslmode=disable", "PostgreSQL datasource name")
	userID := flag.Int("id", 419672615, "Telegram user id")
	flag.Parse()

	return &AppConfig{
		Token:  *token,
		DSN:    *dsn,
		UserID: *userID,
	}
}
