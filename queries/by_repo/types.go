package by_repo

type Query struct {
	shouldBeExecuted bool
	generatedQuery   string
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
