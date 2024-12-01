package aggregator

import (
	"maps"

	"git_applet/gitter"
	"git_applet/queries/mock"
)

func (qa QueryAggregator) GetAll(client gitter.GraphQLClient) (map[string][]gitter.PullRequest, error) {
	prMap := make(map[string][]gitter.PullRequest)

	// mock override
	if qa.Mock {
		prtmp, err := mock.MockQuery{}.GetAll(client)
		if err != nil {
			qa.Logger.Printf("When polling for PRs: %v", err)
		}

		return prtmp, nil
	}

	for _, query := range qa.Queries {

		prtmp, err := query.GetAll(client)
		if err != nil {
			qa.Logger.Printf("When polling for PRs: %v", err)
		}

		maps.Copy(prMap, prtmp)

	}

	return prMap, nil
}
