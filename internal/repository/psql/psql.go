package psql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"musicAPI/internal/config"
	"musicAPI/internal/models"
)

type Storage struct {
	Db *sqlx.DB
}

func MustNewDB(cfg *config.Config) *Storage {
	op := "psql.NewDB"

	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbSSLMode))
	if err != nil {
		t := op + ":" + err.Error()
		panic(t)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		t := op + ":" + err.Error()
		panic(t)
	}

	return &Storage{Db: db}
}

func (s *Storage) AddNewSong(title *models.Title, release string, couplets []string, link string) error {
	// Implementation here
	return nil
}

func (s *Storage) DeleteSong(title *models.Title) error {
	// Implementation here
	return nil
}

func (s *Storage) EditSong(title *models.Title, release string, couplets []string, link string) error {
	// Implementation here
	return nil
}

func (s *Storage) GetCouplets(title *models.Title, page int, limit int) ([]string, error) {
	// Implementation here
	return nil, nil
}

func (s *Storage) GetSongsByGroupsAndRelease(filters *models.Filter, page int, limit int) ([]models.Song, error) {
	// Implementation here
	return nil, nil
}
