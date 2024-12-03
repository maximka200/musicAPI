package psql_test

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"musicAPI/internal/models"
	"musicAPI/internal/repository/psql"
	"testing"
)

func TestGetSongsByGroupsAndRelease(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if sqlxDB == nil {
		t.Fatal("failed to create sqlxDB")
	}
	storage := &psql.Storage{Db: sqlxDB}
	if storage == nil {
		t.Fatal("failed to create storage")
	}

	ctx := context.Background()
	filters := &models.Filter{
		Groups: []string{"group1", "group2"},
		Per: models.Period{
			Start: "1.1.2020",
			End:   "1.1.2020",
		},
	}
	page := 1
	limit := 10

	t.Run("returns songs matching filters", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"group_name", "song_name", "release_date", "link"}).
			AddRow("group1", "song1", "2020-05-01", "link1").
			AddRow("group2", "song2", "2020-06-01", "link2")

		mock.ExpectPrepare("SELECT \\* FROM songs WHERE").
			ExpectQuery().
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), limit, 0).
			WillReturnRows(rows)

		songs, err := storage.GetSongsByGroupsAndRelease(ctx, filters, page, limit)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if len(songs) != 2 {
			t.Fatalf("expected 2 songs, got %d", len(songs))
		}
	})

	t.Run("returns empty slice when no songs match filters", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"group_name", "song_name", "release_date", "link"})

		mock.ExpectPrepare("SELECT \\* FROM songs WHERE").
			ExpectQuery().
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), limit, 0).
			WillReturnRows(rows)

		songs, err := storage.GetSongsByGroupsAndRelease(ctx, filters, page, limit)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if len(songs) != 0 {
			t.Fatalf("expected 0 songs, got %d", len(songs))
		}
	})

	t.Run("returns error on query failure", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM songs WHERE").
			ExpectQuery().
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), limit, 0).
			WillReturnError(fmt.Errorf("query error"))

		_, err := storage.GetSongsByGroupsAndRelease(ctx, filters, page, limit)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
