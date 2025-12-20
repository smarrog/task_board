package main

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/config"
	"github.com/smarrog/task-board/shared/logger"
)

func main() {
	cfg := config.Load()
	dsn := cfg.PostgresDSN
	log := logger.New("core-migrate", zerolog.DebugLevel)

	if dsn == "" {
		log.Fatal().Msg("POSTGRES_DSN is empty")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Error on closing database connection")
		}
	}(db)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal().Err(err).Msg("Error on setting postgres dialect")
	}

	dir := "migrations"
	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "up":
		if err := goose.Up(db, dir); err != nil {
			log.Fatal().Err(err).Msg("Migration failed")
		}
	case "down":
		if err := goose.Down(db, dir); err != nil {
			log.Fatal().Err(err).Msg("Migration failed")
		}
	default:
		log.Fatal().Str("command", cmd).Msg("Unknown command (use up/down)")
	}

	log.Info().Str("command", cmd).Msg("migrations done")
}
