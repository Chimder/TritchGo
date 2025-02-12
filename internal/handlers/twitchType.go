package handlers

type Game struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtURL string `json:"box_art_url"`
	IgdbID    string `json:"igdb_id"`
}

type Pagination struct {
	Cursor string `json:"cursor"`
}

type TopGamesResponse struct {
	Data       []Game     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Stream struct {
	ID           string   `json:"id"`
	UserID       string   `json:"user_id"`
	UserLogin    string   `json:"user_login"`
	UserName     string   `json:"user_name"`
	GameID       string   `json:"game_id"`
	GameName     string   `json:"game_name"`
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	ViewerCount  int      `json:"viewer_count"`
	StartedAt    string   `json:"started_at"`
	Language     string   `json:"language"`
	ThumbnailURL string   `json:"thumbnail_url"`
	TagIDs       []string `json:"tag_ids"`
	Tags         []string `json:"tags"`
	IsMature     bool     `json:"is_mature"`
}

type StreamsResponse struct {
	Data       []Stream   `json:"data"`
	Pagination Pagination `json:"pagination"`
}
