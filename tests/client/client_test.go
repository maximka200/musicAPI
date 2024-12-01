package client_test

import (
	"github.com/stretchr/testify/assert"
	"musicAPI/internal/models"
	"musicAPI/internal/transport/client/musicInfo"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMusicInfo_GetInfo_Success(t *testing.T) {
	expected := &models.Info{
		ReleaseDate: "16.07.2006",
		Text: "Ooh baby, don't you know I suffer?\\nOoh" +
			"baby, can you hear me moan?\\nYou caught me under false pretenses\\n" +
			"How long before you let me go?\\n\\nOoh\\nYou set my soul alight\\" +
			"nOoh\\nYou set my soul alight",
		Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	addr := "http://localhost:1808"

	mi := musicInfo.NewMusicInfo(addr, 2*time.Second)
	title := &models.Title{Group: "TestGroup", Song: "TestSong"}

	info, err := mi.GetInfo(title)

	assert.NoError(t, err)
	assert.Equal(t, expected, info)
}

func TestMusicInfo_GetInfo_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	mi := musicInfo.NewMusicInfo(server.URL, 2*time.Second)
	title := &models.Title{Group: "TestGroup", Song: "TestSong"}

	song, err := mi.GetInfo(title)

	assert.Error(t, err)
	assert.Equal(t, &models.Song{}, song)
}
