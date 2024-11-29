package author

type Query struct {
	generatedQuery string
}

type Config struct {
	Trackers       []Tracker
	ReviewAmount   int
	PrAmount       int
	CommentsAmount int
}

type Tracker struct {
	Id      string
	Authors []string
	Repo    string
	Owner   string
	Title   string
}
