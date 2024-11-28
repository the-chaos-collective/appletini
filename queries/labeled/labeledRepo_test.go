package labeled_test

import (
	"errors"
	"git_applet/queries/labeled"
	"testing"
)

type testConfig struct {
	name          string
	inputs        labeled.Config
	expectedError error
}

func TestValidation(t *testing.T) {
	testCases := []testConfig{
		{
			name: "CommentsAmount > 0",
			inputs: labeled.Config{
				Trackers:       []labeled.Tracker{},
				ReviewAmount:   10,
				PrAmount:       10,
				CommentsAmount: 0,
			},
			expectedError: errors.New("invalid config: CommentsAmount must not be zero or less"),
		},
		{
			name: "PrAmount > 0",
			inputs: labeled.Config{
				Trackers:       []labeled.Tracker{},
				ReviewAmount:   10,
				PrAmount:       0,
				CommentsAmount: 10,
			},
			expectedError: errors.New("invalid config: PrAmount must not be zero or less"),
		},
		{
			name: "ReviewAmount > 0",
			inputs: labeled.Config{
				Trackers:       []labeled.Tracker{},
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
			_, err := labeled.MakeLabeledQuery(testCase.inputs)
			if err.Error() != testCase.expectedError.Error() {
				t.Fatalf("\ngot: %v\nexpected: %v", err.Error(), testCase.expectedError.Error())
			}
		})
	}
}
