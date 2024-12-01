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

	Mergeable string         `yaml:"mergeable"`
	_         map[string]any `yaml:",inline"`
}

func (pr PullRequest) ReviewState() status.ReviewState {
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

func (pr PullRequest) MergeableState() status.MergeableState {
	switch pr.Mergeable {
	case "MERGEABLE":
		return status.MergeableState_Mergeable
	case "CONFLICTING":
		return status.MergeableState_Conflict
	}
	log.Printf("Missing mergeable state: %s", pr.ReviewDecision)
	return status.MergeableState_Unknown
}
