package client_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"musicAPI/internal/models"
	"musicAPI/internal/transport/client/musicInfo"
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

	info, err := mi.GetInfo(context.Background(), title)

	assert.NoError(t, err)
	assert.Equal(t, expected, info)
}

/* SELECT group_name, song_name, release_date, text, link
FROM songs
WHERE
('{"T"}'::text[] IS NULL OR group_name = ANY('{"T"}'::text[]))
AND ('2000-01-01'::date IS NULL OR release_date >= '2000-01-01')
AND ('2024-12-01'::date IS NULL OR release_date <= '2024-12-01')
ORDER BY release_date DESC
LIMIT 10 OFFSET 0; */
