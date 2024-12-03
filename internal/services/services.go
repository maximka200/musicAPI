package services

import (
	"context"
	"fmt"
	"musicAPI/internal/libs/parsers"
	"musicAPI/internal/models"
)

type Repository interface {
	AddNewSong(ctx context.Context, title *models.Title, release string, couplets []string, link string) error
	DeleteSong(ctx context.Context, title *models.Title) error
	EditSong(ctx context.Context, title *models.Title, release string, couplets []string, link string) error
	GetCouplets(ctx context.Context, title *models.Title, page int, limit int) ([]string, error)
	GetSongsByGroupsAndRelease(ctx context.Context, filters *models.Filter, page int, limit int) ([]models.Song, error)
}

type Client interface {
	GetInfo(ctx context.Context, title *models.Title) (*models.Info, error)
}

type Service struct {
	Client Client
	Repos  Repository
}

func NewService(client Client, repos Repository) *Service {
	return &Service{Client: client, Repos: repos}
}

func (s *Service) AddNew(ctx context.Context, title *models.Title) error {
	const op = "service.AddNew"

	info, err := s.Client.GetInfo(ctx, title)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	couplets := parsers.ParseInCouplets(info.Text)

	err = s.Repos.AddNewSong(ctx, title, info.ReleaseDate, couplets, info.Link)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, title *models.Title) error {
	const op = "service.Delete"

	err := s.Repos.DeleteSong(ctx, title)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *Service) Edit(ctx context.Context, song *models.Song) error {
	const op = "service.Edit"
	couplets := parsers.ParseInCouplets(song.Info.Text)
	err := s.Repos.EditSong(ctx, song.Title, song.Info.ReleaseDate, couplets, song.Info.Link)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *Service) GetCouplets(ctx context.Context, title *models.Title, page int, limit int) (string, error) {
	const op = "service.GetCouplets"

	couplets, err := s.Repos.GetCouplets(ctx, title, page, limit)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}

	return parsers.JoinCouplets(couplets), nil
}

func (s *Service) GetSongsByGroupsAndRelease(ctx context.Context, filters *models.Filter, page int, limit int) ([]models.Song, error) {
	const op = "service.GetSongsByGroupsAndRelease"

	filters.Per.Start, _ = parsers.StringDateForPsql(filters.Per.Start)
	filters.Per.End, _ = parsers.StringDateForPsql(filters.Per.End)

	songs, err := s.Repos.GetSongsByGroupsAndRelease(ctx, filters, page, limit)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return songs, nil
}
