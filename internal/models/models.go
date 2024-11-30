package models

import "time"

// inclusive of the end
type Period struct {
	start time.Time
	end   time.Time
}

type Song struct {
	Title Title
	Info  Info
}

type Title struct {
	GroupName string
	SongName  string
}

type Info struct {
	ReleaseDate time.Time
	Couplets    []string
	Link        string
}
