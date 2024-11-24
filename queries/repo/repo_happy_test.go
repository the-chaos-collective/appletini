package repo

import (
	"os"
	"testing"
)

type testCase struct {
	name           string
	inputs         Config
	expectedError  error
	expectedOutput RepoQuery
}

func TestMain(t *testing.M) {
	os.Chdir("../..")

	t.Run()
}

func TestQueryCreation(t *testing.T) {
	testCases := []testCase{
		{
			name:          "happy path",
			expectedError: nil,
			inputs: Config{
				Trackers: []Tracker{
					{
						Name:       "git_appletini",
						Owner:      "darvoid",
						Identifier: "cenas",
					},
				},
				ReviewAmount:   10,
				PrAmount:       10,
				CommentsAmount: 10,
			},
			expectedOutput: RepoQuery{},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			_, err := MakeRepoQuery(test.inputs)
			// if prList == RepoQuery("") {
			// 	t.Fatalf("%v\n#####\nDoes not match:\n%v", prList, test.expectedOutput)
			// }
			if err != test.expectedError {
				t.Fatal("Deu Merda")
			}
		})
	}

}
