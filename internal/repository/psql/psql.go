package psql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"musicAPI/internal/config"
	localError "musicAPI/internal/err"
	"musicAPI/internal/libs/parsers"
	"musicAPI/internal/models"
)

const songsTable = "songs"

type Storage struct {
	db *sqlx.DB
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

	return &Storage{db: db}
}

func (s *Storage) AddNewSong(ctx context.Context, title *models.Title, release string, couplets []string, link string) error {
	const op = "psql.AddNewSong"
	stmt, err := s.db.Prepare(fmt.Sprintf(
		"INSERT INTO %s (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)",
		songsTable))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	date, err := parsers.StringDateForPsql(release)
	if err != nil {
		return err
	}
	if _, err = stmt.ExecContext(
		ctx, title.Group, title.Song, date, pq.Array(couplets), link,
	); err != nil {
		if isUniqueViolation(err) {
			return localError.ErrAlreadyExist
		}

		return fmt.Errorf("%s: %s", op, err.Error())
	}

	return nil
}

func (s *Storage) DeleteSong(ctx context.Context, title *models.Title) error {
	const op = "psql.DeleteSong"
	return nil
}

func (s *Storage) EditSong(ctx context.Context, title *models.Title, release string, couplets []string, link string) error {
	const op = "psql.EditSong"
	return nil
}

func (s *Storage) GetCouplets(ctx context.Context, title *models.Title, page int, limit int) ([]string, error) {
	const op = "psql.GetCouplets"
	return nil, nil
}

func (s *Storage) GetSongsByGroupsAndRelease(ctx context.Context, filters *models.Filter, page int, limit int) ([]models.Song, error) {
	const op = "GetSongsByGroupsAndRelease"

	return nil, nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505" // Код ошибки уникальности в PostgreSQL
	}
	return false
}
