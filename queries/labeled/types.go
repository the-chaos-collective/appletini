package labeled

type LabeledQuery struct {
	generatedQuery string
}

type Config struct {
	Trackers       []Tracker
	ReviewAmount   int
	PrAmount       int
	CommentsAmount int
}

type Tracker struct {
	Label      string
	Repo       string
	Owner      string
	Identifier string
}
