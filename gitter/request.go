package gitter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/machinebox/graphql"
)

// DefaultTimeout bounds a single GraphQL request when no timeout is configured.
const DefaultTimeout = 30 * time.Second

// maxErrorBodyBytes caps how much of a failed response body is kept in the error.
const maxErrorBodyBytes = 512

var (
	ErrEmptyURL         = errors.New("graphql url must not be empty")
	ErrEmptyToken       = errors.New("graphql token must not be empty")
	ErrNegativeTimeout  = errors.New("graphql timeout must not be negative")
	ErrNilResponse      = errors.New("response target must not be nil")
	ErrUnexpectedStatus = errors.New("unexpected HTTP status")
)

// Querier runs a GraphQL query and decodes the "data" field into response.
// Queries depend on this interface so they can be tested without a network.
type Querier interface {
	Query(ctx context.Context, query string, response any) error
}

// GraphQLClient is a Querier backed by a single, reused HTTP client.
type GraphQLClient struct {
	client  *graphql.Client
	token   string
	timeout time.Duration
}

// NewGraphQLClient builds a client for url authenticated with token.
// Each request is bounded by timeout; zero means DefaultTimeout.
func NewGraphQLClient(url, token string, timeout time.Duration) (*GraphQLClient, error) {
	if url == "" {
		return nil, ErrEmptyURL
	}

	if token == "" {
		return nil, ErrEmptyToken
	}

	if timeout < 0 {
		return nil, ErrNegativeTimeout
	}

	if timeout == 0 {
		timeout = DefaultTimeout
	}

	// No http.Client.Timeout: each request gets a context deadline instead,
	// so callers can also cancel it.
	//exhaustruct:ignore
	httpClient := &http.Client{
		Transport: statusCheckingTransport{next: http.DefaultTransport},
	}

	return &GraphQLClient{
		client:  graphql.NewClient(url, graphql.WithHTTPClient(httpClient)),
		token:   token,
		timeout: timeout,
	}, nil
}

// Query runs query, honouring ctx cancellation and the client's per-request timeout.
func (c *GraphQLClient) Query(ctx context.Context, query string, response any) error {
	if response == nil {
		return ErrNilResponse
	}

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("before executing GraphQL query: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req := graphql.NewRequest(query)
	req.Header.Set("Authorization", "Bearer "+c.token)

	err = c.client.Run(ctx, req, response)
	if err != nil {
		return fmt.Errorf("executing GraphQL query: %w", err)
	}

	return nil
}

// statusCheckingTransport turns non-2xx responses into errors. The graphql
// library ignores the status code, so without this a 401 (bad token) decodes
// as an empty, successful result.
type statusCheckingTransport struct {
	next http.RoundTripper
}

func (t statusCheckingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.next.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return res, nil
	}

	defer res.Body.Close()

	body, readErr := io.ReadAll(io.LimitReader(res.Body, maxErrorBodyBytes))
	if readErr != nil {
		return nil, fmt.Errorf("%w: %s (reading body: %w)", ErrUnexpectedStatus, res.Status, readErr)
	}

	return nil, fmt.Errorf("%w: %s: %s", ErrUnexpectedStatus, res.Status, body)
}
