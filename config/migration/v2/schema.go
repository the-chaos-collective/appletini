package v2

type Config struct {
	Computed  Computed     `json:"-"`
	Version   int          `json:"__version"`
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
	Personal bool       `json:"personal"`
	ByLabel  LabeledSet `json:"byLabel"`
	ByRepo   RepoSet    `json:"byRepo"`
	ByAuthor AuthorSet  `json:"byAuthor"`
}

type Author struct {
	Title    string   `json:"title"`
	Owner    string   `json:"owner"`
	RepoName string   `json:"repo"`
	Authors  []string `json:"authors"`
}

type Labeled struct {
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

type (
	LabeledSet []Labeled
	RepoSet    []Repo
	AuthorSet  []Author
)

type Computed struct {
	GithubToken string
}
