package types

type Config struct {
	Github    GithubConfig `json:"github"`
	Poll      PollConfig   `json:"poll"`
	Tracking  Tracking     `json:"tracking"`
	ItemCount int          `json:"item_count"`
	Darkmode  bool         `json:"darkmode"`
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
	ByLabel        LabeledRepoSet `json:"byLabel"`
	ByRepo         RepoSet        `json:"byRepo"`
	Projects       ProjectSet     `json:"projects"`
	CommentsAmount string         `json:"commentsAmount"`
	ReviewAmount   string         `json:"reviewAmount"`
	PrAmount       string         `json:"prAmount"`
}
type LabeledRepo struct {
	Identifier string `json:"identifier"`
	Owner      string `json:"owner"`
	RepoName   string `json:"repo"`
	Label      string `json:"label"`
}
type Repo struct {
	Identifier string `json:"identifier"`
	Owner      string `json:"owner"`
	RepoName   string `json:"repo"`
}

type Project struct {
	Identifier string `json:"identifier"`
	Owner      string `json:"owner"`
	RepoName   string `json:"repo"`
}

type (
	LabeledRepoSet []LabeledRepo
	RepoSet        []Repo
	ProjectSet     []Project
)
