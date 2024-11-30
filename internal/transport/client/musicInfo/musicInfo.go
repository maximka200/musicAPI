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

func (mi *MusicInfo) GetInfo(title models.Title) (models.Song, error) {
	const op = "musicInfo.GetInfo"

	url := fmt.Sprintf("%s/info?group=%s&song=%s", mi.Address, title.Group, title.Song)
	resp, err := mi.Client.Get(url)
	if err != nil {
		return models.Song{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Song{}, fmt.Errorf("failed to get info: %s", resp.Status)
	}

	var song models.Song
	if err = json.NewDecoder(resp.Body).Decode(&song); err != nil {
		return models.Song{}, err
	}

	return song, nil
}
