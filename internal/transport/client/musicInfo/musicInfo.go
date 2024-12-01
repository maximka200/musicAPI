package musicInfo

import (
	"encoding/json"
	"fmt"
	"musicAPI/internal/models"
	"net/http"
	"time"
)

type MusicInfo struct {
	Client  *http.Client
	Address string
}

func NewMusicInfo(address string, timeout time.Duration) *MusicInfo {
	client := &http.Client{Timeout: timeout}
	return &MusicInfo{
		Client:  client,
		Address: address,
	}
}

func (mi *MusicInfo) GetInfo(title *models.Title) (*models.Song, error) {
	const op = "musicInfo.GetInfo"

	url := fmt.Sprintf("%s/info?group=%s&song=%s", mi.Address, title.Group, title.Song)
	resp, err := mi.Client.Get(url)
	if err != nil {
		return &models.Song{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &models.Song{}, fmt.Errorf("%s: %w", op, err)
	}

	info := &models.Info{}
	if err = json.NewDecoder(resp.Body).Decode(info); err != nil {
		return &models.Song{}, fmt.Errorf("%s: %w", op, err)
	}

	return &models.Song{Title: title, Info: info}, nil
}
