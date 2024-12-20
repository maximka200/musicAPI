package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"musicAPI/internal/config"
	localError "musicAPI/internal/err"
	"musicAPI/internal/libs/parsers"
	"musicAPI/internal/models"
)

const songsTable = "songs"

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

func (s *Storage) AddNewSong(ctx context.Context, title *models.Title, release string, couplets []string, link string) error {
	const op = "psql.AddNewSong"
	stmt, err := s.Db.Prepare(fmt.Sprintf(
		"INSERT INTO %s (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)",
		songsTable))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	date, err := parsers.StringDateForPsql(release)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if _, err = stmt.Exec(
		title.Group, title.Song, date, pq.Array(couplets), link,
	); err != nil {
		if isUniqueViolation(err) {
			return localError.ErrAlreadyExist
		}

		return fmt.Errorf("%s: %s %s", op, err.Error())
	}

	return nil
}

func (s *Storage) DeleteSong(ctx context.Context, title *models.Title) error {
	const op = "psql.DeleteSong"
	stmt, err := s.Db.Prepare(fmt.Sprintf(
		"DELETE FROM %s WHERE group_name=$1 AND song_name=$2",
		songsTable))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(
		title.Group, title.Song)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return localError.ErrNotFound
		}
		return fmt.Errorf("%s: %s", op, err.Error())
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %s %s", op, err.Error())
	}
	if rowsAff == 0 {
		return localError.ErrNotFound
	}

	return nil
}

func (s *Storage) EditSong(ctx context.Context, title *models.Title, release string, couplets []string, link string) error {
	const op = "psql.EditSong"
	stmt, err := s.Db.Prepare(fmt.Sprintf(
		"UPDATE %s SET release_date = $1, text = $2, link = $3 WHERE group_name = $4 AND song_name = $5",
		songsTable))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	date, err := parsers.StringDateForPsql(release)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	result, err := stmt.Exec(
		date, pq.Array(couplets), link,
		title.Group, title.Song)
	if err != nil {
		return fmt.Errorf("%s: %s %s", op, err.Error())
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %s %s", op, err.Error())
	}
	if rowsAff == 0 {
		return localError.ErrNotFound
	}

	return nil
}

func (s *Storage) GetCouplets(ctx context.Context, title *models.Title, page, limit int) ([]string, error) {
	const op = "psql.GetCouplets"

	offset := (page - 1) * limit
	stmt, err := s.Db.Prepare(fmt.Sprintf(`
		SELECT couplet 
		FROM unnest((
			SELECT text 
			FROM %s
			WHERE group_name = $1 AND song_name = $2
		)) AS couplet
		LIMIT $3 OFFSET $4;
		`,
		songsTable))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(title.Group, title.Song, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var couplets []string
	for rows.Next() {
		var couplet string
		if err = rows.Scan(&couplet); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}
		couplets = append(couplets, couplet)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration failed: %w", op, err)
	}

	if len(couplets) == 0 {
		return nil, localError.ErrNotFound
	}

	return couplets, nil
}

func (s *Storage) GetSongsByGroupsAndRelease(ctx context.Context, filters *models.Filter, page int, limit int) ([]models.Song, error) {
	const op = "psql.GetSongsByGroupsAndRelease"

	offset := (page - 1) * limit

	query := fmt.Sprintf(`
		SELECT group_name, song_name, release_date, text, link
		FROM %s
		WHERE
		($1::text[] IS NULL OR group_name = ANY($1))
		AND (NULLIF($2, '')::date IS NULL OR release_date >= NULLIF($2, '')::date)
		AND (NULLIF($3, '')::date IS NULL OR release_date <= NULLIF($3, '')::date)
		ORDER BY release_date DESC
		LIMIT $4 OFFSET $5;
	`, songsTable)

	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare query: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		pq.Array(filters.Groups),
		filters.Per.Start,
		filters.Per.End,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w, %s, %s", op, err,
			filters.Per.End, filters.Per.Start)
	}
	defer rows.Close()

	songs := []models.Song{}
	for rows.Next() {
		song := models.Song{
			Title: &models.Title{},
			Info:  &models.Info{},
		}
		bufferCouplets := []string{}
		bufferRelease := ""
		if err := rows.Scan(
			&song.Title.Group,
			&song.Title.Song,
			&bufferRelease,
			pq.Array(&bufferCouplets),
			&song.Info.Link,
		); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}

		song.Info.Text = parsers.JoinCouplets(bufferCouplets)
		song.Info.ReleaseDate, err = parsers.ConvertISOToDate(bufferRelease)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return songs, nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}
