package personal

import "appletini/gitter"

type Response struct {
	Viewer struct {
		PullRequests struct {
			Nodes []edge `json:"edges"`
		} `json:"pullRequests"`
	} `json:"viewer"`
}
type edge struct {
	Node pullRequestPersonal `yaml:"node"`
}

type pullRequestPersonal struct {
	Id          string `yaml:"id"`
	Title       string `yaml:"title"`
	BaseRefName string `yaml:"baseRefName"`
	HeadRefName string `yaml:"headRefName"`
	Number      int    `yaml:"number"`
	Permalink   string `yaml:"permalink"`
	ReviewCount struct {
		TotalCount int `yaml:"totalCount"`
	} `yaml:"reviewCount"`
	ReviewRequests struct {
		TotalCount int `yaml:"totalCount"`
	} `yaml:"reviewRequests"`
	ReviewDecision string `yaml:"reviewDecision"`
	Mergeable      string `yaml:"mergeable"`
}

func (pr pullRequestPersonal) transform() gitter.PullRequest {
	return gitter.PullRequest{
		Id:             pr.Id,
		Title:          pr.Title,
		BaseRefName:    pr.BaseRefName,
		HeadRefName:    pr.HeadRefName,
		Number:         pr.Number,
		Permalink:      pr.Permalink,
		ReviewCount:    pr.ReviewCount.TotalCount,
		ReviewRequests: pr.ReviewRequests.TotalCount,
		ReviewDecision: pr.ReviewDecision,
		Mergeable:      pr.Mergeable,
	}
}

func (pr Response) Extract() []gitter.PullRequest {
	prs := []gitter.PullRequest{}
	for _, val := range pr.Viewer.PullRequests.Nodes {
		prs = append(prs, val.Node.transform())
	}

	return prs
}
