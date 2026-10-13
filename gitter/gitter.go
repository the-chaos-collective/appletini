package gitter

import (
	"log"

	"appletini/status"
)

type PullRequest struct {
	Title       string `yaml:"title"`
	BaseRefName string `yaml:"baseRefName"`
	HeadRefName string `yaml:"headRefName"`
	Number      int    `yaml:"number"`
	Permalink   string `yaml:"permalink"`
	ReviewCount int    `yaml:"reviewCount"`

	ReviewRequests int    `yaml:"reviewRequests"`
	ReviewDecision string `yaml:"reviewDecision"`
	Id             string `yaml:"id"`

	Mergeable string `yaml:"mergeable"`

	// ChecksState is GitHub's combined CI state on the latest commit
	// (SUCCESS, FAILURE, ERROR, PENDING, EXPECTED), or "" when it has no checks.
	ChecksState string `yaml:"checksState"`

	// ViewerReview is the state of your own latest review on the PR
	// (APPROVED, CHANGES_REQUESTED, COMMENTED, DISMISSED, PENDING), or ""
	// when you haven't reviewed it.
	ViewerReview string         `yaml:"viewerReview"`
	_            map[string]any `yaml:",inline"`
}

// Review is the response shape of `viewerLatestReview { state }`.
type Review struct {
	State string `json:"state"`
}

// StateOrEmpty returns the review's state, or "" when there is no review.
func (r *Review) StateOrEmpty() string {
	if r == nil {
		return ""
	}

	return r.State
}

// ApprovedByViewer reports whether your own latest review approved the PR.
// Other people's approvals don't count.
func (pr PullRequest) ApprovedByViewer() bool {
	return pr.ViewerReview == "APPROVED"
}

// ApprovableByViewer reports whether offering you an Approve makes sense:
// not if you've already approved, and not while the PR is red because of
// merge conflicts, failing CI or changes requested by someone else. Changes
// you requested yourself don't block it, since GitHub keeps the PR in
// "changes requested" until you approve (or your review is dismissed).
func (pr PullRequest) ApprovableByViewer() bool {
	if pr.ApprovedByViewer() {
		return false
	}

	info := pr.PRInfo()

	if info.Mergeable == status.MergeableState_Conflict || info.Checks == status.ChecksState_Failing {
		return false
	}

	if info.Review == status.ReviewState_ChangesRequested && pr.ViewerReview != "CHANGES_REQUESTED" {
		return false
	}

	return true
}

// LastCommit is the response shape of
//
//	commits(last: 1) { nodes { commit { statusCheckRollup { state } } } }
//
// which every PR query requests to get the PR's CI state.
type LastCommit struct {
	Nodes []struct {
		Commit struct {
			StatusCheckRollup *struct {
				State string `json:"state"`
			} `json:"statusCheckRollup"`
		} `json:"commit"`
	} `json:"nodes"`
}

// ChecksState returns the combined CI state of the latest commit, or "" when
// there is no commit or no checks ran on it.
func (c LastCommit) ChecksState() string {
	if len(c.Nodes) == 0 || c.Nodes[0].Commit.StatusCheckRollup == nil {
		return ""
	}

	return c.Nodes[0].Commit.StatusCheckRollup.State
}

func (pr PullRequest) reviewState() status.ReviewState {
	switch pr.ReviewDecision {
	case "APPROVED":
		return status.ReviewState_Approved
	case "REVIEW_REQUIRED":
		return status.ReviewState_RequiresReview
	case "CHANGES_REQUESTED":
		return status.ReviewState_ChangesRequested
	case "":
		return status.ReviewState_NoReviewRequired
	}
	log.Printf("Missing review state: %s", pr.ReviewDecision)
	return status.ReviewState_Unknown
}

func (pr PullRequest) mergeableState() status.MergeableState {
	switch pr.Mergeable {
	case "MERGEABLE":
		return status.MergeableState_Mergeable
	case "CONFLICTING":
		return status.MergeableState_Conflict
	case "UNKNOWN":
		return status.MergeableState_Unknown
	}
	log.Printf("Missing mergeable state: %s", pr.Mergeable)
	return status.MergeableState_Unknown
}

func (pr PullRequest) checksState() status.ChecksState {
	switch pr.ChecksState {
	case "":
		return status.ChecksState_None
	case "SUCCESS":
		return status.ChecksState_Passing
	case "FAILURE", "ERROR":
		return status.ChecksState_Failing
	case "PENDING", "EXPECTED":
		return status.ChecksState_Pending
	}
	log.Printf("Missing checks state: %s", pr.ChecksState)
	return status.ChecksState_Unknown
}

func (pr PullRequest) PRInfo() status.PRInfo {
	return status.PRInfo{
		Review:    pr.reviewState(),
		Mergeable: pr.mergeableState(),
		Checks:    pr.checksState(),
	}
}
