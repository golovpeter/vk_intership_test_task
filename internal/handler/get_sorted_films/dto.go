package get_sorted_films

type FilmData struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ReleaseDate string   `json:"release_date"`
	Rating      int      `json:"rating"`
	Actors      []string `json:"actors"`
}
