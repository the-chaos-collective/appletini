package repo

type RepoQuery struct {
	generatedQuery string
}

type Config struct {
	Trackers       []Tracker
	ReviewAmount   int
	PrAmount       int
	CommentsAmount int
}

type Tracker struct {
	Id    string
	Repo  string
	Owner string
	Title string
}
