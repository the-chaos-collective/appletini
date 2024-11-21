package types

type Config struct {
	ChromeProfile string       `json:"chrome_profile"`
	Github        GithubConfig `json:"github"`
	Poll          PollConfig   `json:"poll"`
	Tracking      Tracking     `json:"tracking"`
	ItemCount     int          `json:"item_count"`
}

type GithubConfig struct {
	Host     string `json:"host"`
	GraphQL  string `json:"gqlAPI"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type PollConfig struct {
	Frequency int `json:"frequency_s"`
}
type Tracking struct {
	ByLabel        RepoSet          `json:"ByLabel"`
	ByRepo         UnlabeledRepoSet `json:"byRepo"`
	Project        ProjectSet       `json:"project"`
	CommentsAmount string           `json:"commentsAmount"`
	ReviewAmount   string           `json:"reviewAmount"`
	PrAmount       string           `json:"prAmount"`
}
type Repo struct {
	Owner      string `json:"owner"`
	RepoName   string `json:"repo"`
	Identifier string `json:"identifier"`
	Label      string `json:"label"`
}
type UnlabeledRepo struct {
	Owner      string `json:"owner"`
	RepoName   string `json:"repo"`
	Identifier string `json:"identifier"`
}

type Project struct {
	Owner      string `json:"owner"`
	RepoName   string `json:"repo"`
	Identifier string `json:"identifier"`
}

type (
	RepoSet          []Repo
	UnlabeledRepoSet []UnlabeledRepo
	ProjectSet       []Project
)
