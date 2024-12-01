package v1

type Config struct {
	Github    GithubConfig `json:"github"`
	Poll      PollConfig   `json:"poll"`
	Tracking  Tracking     `json:"tracking"`
	ItemCount int          `json:"itemCount"`
	Darkmode  bool         `json:"darkMode"`
}

type GithubConfig struct {
	GraphQL string `json:"gqlAPI"`
	Token   string `json:"token"`
}

type PollConfig struct {
	Frequency int `json:"frequencySeconds"`
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
	Title    string `json:"title"`
	Owner    string `json:"owner"`
	RepoName string `json:"repo"`
	Label    string `json:"label"`
}
type Repo struct {
	Title    string `json:"title"`
	Owner    string `json:"owner"`
	RepoName string `json:"repo"`
}

type Project struct {
	Title    string `json:"title"`
	Owner    string `json:"owner"`
	RepoName string `json:"repo"`
}

type (
	LabeledRepoSet []LabeledRepo
	RepoSet        []Repo
	ProjectSet     []Project
)
