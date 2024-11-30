package services

import (
	"musicAPI/internal/models"
	"musicAPI/internal/transport/client/musicInfo"
)

// TODO: parser of couplets from text through "\n\n"

type Service struct {
	client *musicInfo.MusicInfo
	// repos
}

func NewService(client *musicInfo.MusicInfo) *Service {
	return &Service{client: client}
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

func (s *Service) GetCouplet(title *models.Title, startInd int, endInd int) (string, error) {
	panic("not impl")
}

func (s *Service) GetSongsByGroupsAndRelease(groups []string, period models.Period) ([]models.Song, error) {
	panic("not impl")
}
