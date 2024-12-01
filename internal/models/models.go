package models

import "time"

type Filter struct {
	Groups []string
	Per    Period
}

// inclusive of the end
type Period struct {
	Start time.Time
	End   time.Time
}

type Song struct {
	Title *Title `json:"title"`
	Info  *Info  `json:"info"`
}

type Title struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type Info struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
