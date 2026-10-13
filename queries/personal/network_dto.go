package personal

import "appletini/gitter"

type Response struct {
	Viewer struct {
		PullRequests struct {
			Edges []edge `json:"edges"`
		} `json:"pullRequests"`
	} `json:"viewer"`
}

type edge struct {
	Node pullRequestPersonal `json:"node"`
}

type totalCount struct {
	TotalCount int `json:"totalCount"`
}

type pullRequestPersonal struct {
	Id             string            `json:"id"`
	Title          string            `json:"title"`
	BaseRefName    string            `json:"baseRefName"`
	HeadRefName    string            `json:"headRefName"`
	Number         int               `json:"number"`
	Permalink      string            `json:"permalink"`
	Reviews        totalCount        `json:"reviews"`
	ReviewRequests totalCount        `json:"reviewRequests"`
	ReviewDecision string            `json:"reviewDecision"`
	Mergeable      string            `json:"mergeable"`
	Commits        gitter.LastCommit `json:"commits"`
	ViewerReview   *gitter.Review    `json:"viewerLatestReview"`
}

func (pr pullRequestPersonal) transform() gitter.PullRequest {
	return gitter.PullRequest{
		Id:             pr.Id,
		Title:          pr.Title,
		BaseRefName:    pr.BaseRefName,
		HeadRefName:    pr.HeadRefName,
		Number:         pr.Number,
		Permalink:      pr.Permalink,
		ReviewCount:    pr.Reviews.TotalCount,
		ReviewRequests: pr.ReviewRequests.TotalCount,
		ReviewDecision: pr.ReviewDecision,
		Mergeable:      pr.Mergeable,
		ChecksState:    pr.Commits.ChecksState(),
		ViewerReview:   pr.ViewerReview.StateOrEmpty(),
		ViewerAuthor:   true, // viewer.pullRequests only lists your own
	}
}

func (pr Response) Extract() []gitter.PullRequest {
	prs := []gitter.PullRequest{}
	for _, val := range pr.Viewer.PullRequests.Edges {
		prs = append(prs, val.Node.transform())
	}

	return prs
}
