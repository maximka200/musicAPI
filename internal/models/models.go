package models

import "time"

// inclusive of the end
type Period struct {
	start time.Time
	end   time.Time
}

type Song struct {
	Title Title
	Info  interface{} // InfoCoupletsArray or InfoCoupletsString
}

type Title struct {
	Group string
	Song  string
}

type InfoCoupletsArray struct {
	ReleaseDate time.Time
	Couplets    []string
	Link        string
}

type InfoCoupletsString struct {
	ReleaseDate time.Time
	Couplets    string
	Link        string
}
