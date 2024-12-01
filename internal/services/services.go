package services

import (
	"musicAPI/internal/models"
	"time"
)

type Repository interface {
	AddNewSong(title models.Title, release time.Time, couplets []string, link string) error
	DeleteSong(title models.Title) error
	EditSong(title models.Title, release time.Time, couplets []string, link string) error
	GetCouplets(title models.Title, page int, limit int) (string, error)
	GetSongsByGroupsAndRelease(filters *models.Filter, page int, limit int) ([]models.Song, error)
}

type Client interface {
	GetInfo(title models.Title) (*models.Song, error)
}

// TODO: parser of couplets from text through "\n\n"

type Service struct {
	Client Client
	Repos  Repository
}

func NewService(client Client, repos Repository) *Service {
	return &Service{Client: client, Repos: repos}
}

func (s *Service) AddNew(title *models.Title) error {
	panic("not impl")
}

func (s *Service) Delete(title *models.Title) error {
	panic("not impl")
}

func (s *Service) Edit(song *models.Song) error {
	panic("not impl")
}

func (s *Service) GetCouplet(title *models.Title, page int, limit int) (string, error) {
	panic("not impl")
}

func (s *Service) GetSongsByGroupsAndRelease(filters *models.Filter, page int, limit int) ([]models.Song, error) {
	panic("not impl")
}
