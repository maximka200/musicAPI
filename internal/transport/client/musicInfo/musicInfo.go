package musicInfo

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	localError "musicAPI/internal/err"
	"musicAPI/internal/models"
	"net/http"
	"net/url"
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

func (mi *MusicInfo) GetInfo(ctx context.Context, title *models.Title) (*models.Info, error) {
	const op = "musicInfo.GetInfo"
	slog.Info(title.Song, title.Group, mi.Address)
	url := fmt.Sprintf("%s/info?group=%s&song=%s", mi.Address,
		url.QueryEscape(title.Group),
		url.QueryEscape(title.Song))
	resp, err := mi.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf("%s:%w", op, localError.ErrBadRequest)
		}
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%s:%w", op, localError.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	info := &models.Info{}
	if err = json.NewDecoder(resp.Body).Decode(info); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return info, nil
}
