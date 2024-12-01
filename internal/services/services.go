package services

import (
	"fmt"
	"musicAPI/internal/libs/parsers"
	"musicAPI/internal/models"
)

type Repository interface {
	AddNewSong(title *models.Title, release string, couplets []string, link string) error
	DeleteSong(title *models.Title) error
	EditSong(title *models.Title, release string, couplets []string, link string) error
	GetCouplets(title *models.Title, page int, limit int) ([]string, error)
	GetSongsByGroupsAndRelease(filters *models.Filter, page int, limit int) ([]models.Song, error)
}

type Client interface {
	GetInfo(title *models.Title) (*models.Info, error)
}

type Service struct {
	Client Client
	Repos  Repository
}

func NewService(client Client, repos Repository) *Service {
	return &Service{Client: client, Repos: repos}
}

func (s *Service) AddNew(title *models.Title) error {
	const op = "service.AddNew"

	info, err := s.Client.GetInfo(title)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	couplets := parsers.ParseInCouplets(info.Text)

	err = s.Repos.AddNewSong(title, info.ReleaseDate, couplets, info.Link)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *Service) Delete(title *models.Title) error {
	const op = "service.Delete"

	err := s.Repos.DeleteSong(title)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *Service) Edit(song *models.Song) error {
	const op = "service.Edit"
	couplets := parsers.ParseInCouplets(song.Info.Text)
	err := s.Repos.EditSong(song.Title, song.Info.Text, couplets, song.Info.Link)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *Service) GetCouplets(title *models.Title, page int, limit int) (string, error) {
	const op = "service.GetCouplets"

	couplets, err := s.Repos.GetCouplets(title, page, limit)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}

	return parsers.JoinCouplets(couplets), nil
}

func (s *Service) GetSongsByGroupsAndRelease(filters *models.Filter, page int, limit int) ([]models.Song, error) {
	const op = "service.GetSongsByGroupsAndReleaze"

	songs, err := s.Repos.GetSongsByGroupsAndRelease(filters, page, limit)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return songs, nil
}
