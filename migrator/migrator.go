package migrator

import (
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"musicAPI/internal/config"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func MigrateUp(cfg config.Config, vendor string) {
	connectionStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		vendor, cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.DbSSLMode)
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		panic("cannot connect with db:" + err.Error())
	}
	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("cannot ping db: %v", err))
	}

	goose.SetBaseFS(embedMigrations)

	if err = goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err = goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	db.Close()
}

func MigrateDown(cfg config.Config, vendor string) {
	connectionStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		vendor, cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.DbSSLMode)
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		panic("cannot connect with db")
	}
	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("cannot ping db: %v", err))
	}

	goose.SetBaseFS(embedMigrations)

	if err = goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err = goose.Down(db, "migrations"); err != nil {
		panic(err)
	}

	db.Close()
}
