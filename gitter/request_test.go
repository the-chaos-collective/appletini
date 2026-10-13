package gitter_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"appletini/gitter"
)

func TestNewGraphQLClient(t *testing.T) {
	testCases := []struct {
		name        string
		url         string
		token       string
		timeout     time.Duration
		expectedErr error
	}{
		{name: "valid", url: "https://api.github.com/graphql", token: "t", timeout: time.Second},
		{name: "zero timeout uses default", url: "https://api.github.com/graphql", token: "t"},
		{name: "empty url", url: "", token: "t", expectedErr: gitter.ErrEmptyURL},
		{name: "empty token", url: "https://x", token: "", expectedErr: gitter.ErrEmptyToken},
		{name: "negative timeout", url: "https://x", token: "t", timeout: -time.Second, expectedErr: gitter.ErrNegativeTimeout},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := gitter.NewGraphQLClient(tc.url, tc.token, tc.timeout)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if (client == nil) == (tc.expectedErr == nil) {
				t.Fatalf("client nil=%v with error %v", client == nil, err)
			}
		})
	}
}

type viewer struct {
	Viewer struct {
		Login string `json:"login"`
	} `json:"viewer"`
}

func TestGraphQLClientQuery(t *testing.T) {
	const token = "secret-token"

	// handler factory: checks the request and answers with status/body after delay.
	respond := func(status int, body string, delay time.Duration) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method %s, expected POST", r.Method)
			}

			if got := r.Header.Get("Authorization"); got != "Bearer "+token {
				t.Errorf("authorization header %q", got)
			}

			reqBody, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(reqBody), "viewer") {
				t.Errorf("request body %s does not contain the query", reqBody)
			}

			if delay > 0 {
				select {
				case <-time.After(delay):
				case <-r.Context().Done():
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = io.WriteString(w, body)
		}
	}

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()

	testCases := []struct {
		name          string
		handler       http.HandlerFunc
		ctx           context.Context
		timeout       time.Duration
		nilResponse   bool
		expectedLogin string
		expectedErr   error
		errContains   string
	}{
		{
			name:          "success",
			handler:       respond(http.StatusOK, `{"data":{"viewer":{"login":"octocat"}}}`, 0),
			ctx:           context.Background(),
			expectedLogin: "octocat",
		},
		{
			name:        "graphql error",
			handler:     respond(http.StatusOK, `{"errors":[{"message":"Could not resolve to a Repository"}]}`, 0),
			ctx:         context.Background(),
			errContains: "Could not resolve to a Repository",
		},
		{
			name:        "bad credentials (401 with JSON body)",
			handler:     respond(http.StatusUnauthorized, `{"message":"Bad credentials"}`, 0),
			ctx:         context.Background(),
			expectedErr: gitter.ErrUnexpectedStatus,
			errContains: "Bad credentials",
		},
		{
			name:        "server error",
			handler:     respond(http.StatusBadGateway, `bad gateway`, 0),
			ctx:         context.Background(),
			expectedErr: gitter.ErrUnexpectedStatus,
		},
		{
			name:        "invalid json",
			handler:     respond(http.StatusOK, `not json`, 0),
			ctx:         context.Background(),
			errContains: "executing GraphQL query",
		},
		{
			name:        "already cancelled",
			handler:     respond(http.StatusOK, `{"data":{}}`, 0),
			ctx:         cancelled,
			expectedErr: context.Canceled,
		},
		{
			name:        "timeout",
			handler:     respond(http.StatusOK, `{"data":{}}`, 2*time.Second),
			ctx:         context.Background(),
			timeout:     20 * time.Millisecond,
			expectedErr: context.DeadlineExceeded,
		},
		{
			name:        "nil response",
			handler:     respond(http.StatusOK, `{"data":{}}`, 0),
			ctx:         context.Background(),
			nilResponse: true,
			expectedErr: gitter.ErrNilResponse,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client, err := gitter.NewGraphQLClient(server.URL, token, tc.timeout)
			if err != nil {
				t.Fatalf("creating client: %v", err)
			}

			var res viewer

			var target any = &res
			if tc.nilResponse {
				target = nil
			}

			err = client.Query(tc.ctx, `query { viewer { login } }`, target)

			if tc.expectedErr != nil && !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if tc.errContains != "" && (err == nil || !strings.Contains(err.Error(), tc.errContains)) {
				t.Fatalf("got error %v, expected it to contain %q", err, tc.errContains)
			}

			if tc.expectedErr == nil && tc.errContains == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res.Viewer.Login != tc.expectedLogin {
				t.Fatalf("got login %q, expected %q", res.Viewer.Login, tc.expectedLogin)
			}
		})
	}
}
