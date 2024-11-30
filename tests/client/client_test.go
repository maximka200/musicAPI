package client_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"musicAPI/internal/models"
	"musicAPI/internal/transport/client/musicInfo"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMusicInfo_GetInfo_Success(t *testing.T) {
	expectedSong := models.Song{Title: models.Title{Group: "TestGroup", Song: "TestSong"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/info?group=TestGroup&song=TestSong", r.URL.String())
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedSong)
	}))
	defer server.Close()

	mi := musicInfo.NewMusicInfo(server.URL, 2*time.Second)
	title := models.Title{Group: "TestGroup", Song: "TestSong"}

	song, err := mi.GetInfo(title)

	assert.NoError(t, err)
	assert.Equal(t, expectedSong, song)
}

func TestMusicInfo_GetInfo_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	mi := musicInfo.NewMusicInfo(server.URL, 2*time.Second)
	title := models.Title{Group: "TestGroup", Song: "TestSong"}

	song, err := mi.GetInfo(title)

	assert.Error(t, err)
	assert.Equal(t, models.Song{}, song)
}

func TestMusicInfo_GetInfo_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	mi := musicInfo.NewMusicInfo(server.URL, 2*time.Second)
	title := models.Title{Group: "TestGroup", Song: "TestSong"}

	song, err := mi.GetInfo(title)

	assert.Error(t, err)
	assert.Equal(t, models.Song{}, song)
}

func TestMusicInfo_GetInfo_HTTPError(t *testing.T) {
	mi := musicInfo.NewMusicInfo("http://invalid-url", 2*time.Second)
	title := models.Title{Group: "TestGroup", Song: "TestSong"}

	song, err := mi.GetInfo(title)

	assert.Error(t, err)
	assert.Equal(t, models.Song{}, song)
}
