package models

type Filter struct {
	Groups []string `json:"groups"`
	Per    *Period  `json:"period"`
}

// inclusive of the end
type Period struct {
	Start string `json:"start"`
	End   string `json:"end"`
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
