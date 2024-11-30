package author_test

import (
	"errors"
	"testing"

	"git_applet/queries/author"
)

type testConfig struct {
	name          string
	inputs        author.Config
	expectedError error
}

func TestValidation(t *testing.T) {
	testCases := []testConfig{
		{
			name: "Trackers.Authors > 0",
			inputs: author.Config{
				Trackers: []author.Tracker{
					{
						Id:      "foo",
						Authors: []string{},
						Repo:    "foo",
						Owner:   "foo",
						Title:   "foo",
					},
				},
				ReviewAmount:   10,
				PrAmount:       10,
				CommentsAmount: 10,
			},
			expectedError: errors.New("invalid config: Trackers[0].Authors must have at least one author"),
		},
		{
			name: "CommentsAmount > 0",
			inputs: author.Config{
				Trackers:       []author.Tracker{},
				ReviewAmount:   10,
				PrAmount:       10,
				CommentsAmount: 0,
			},
			expectedError: errors.New("invalid config: CommentsAmount must not be zero or less"),
		},
		{
			name: "PrAmount > 0",
			inputs: author.Config{
				Trackers:       []author.Tracker{},
				ReviewAmount:   10,
				PrAmount:       0,
				CommentsAmount: 10,
			},
			expectedError: errors.New("invalid config: PrAmount must not be zero or less"),
		},
		{
			name: "ReviewAmount > 0",
			inputs: author.Config{
				Trackers:       []author.Tracker{},
				ReviewAmount:   0,
				PrAmount:       10,
				CommentsAmount: 10,
			},
			expectedError: errors.New("invalid config: ReviewAmount must not be zero or less"),
		},
		// TODO: add more cases
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := author.MakeQuery(testCase.inputs)
			if err.Error() != testCase.expectedError.Error() {
				t.Fatalf("\ngot: %v\nexpected: %v", err.Error(), testCase.expectedError.Error())
			}
		})
	}
}
