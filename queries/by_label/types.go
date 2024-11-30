package by_label

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
	Label string
	Repo  string
	Owner string
	Title string
}
