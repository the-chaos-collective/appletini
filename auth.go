package main

import (
	"context"
	"os"

	"golang.org/x/oauth2"
)

func auth2() (ctx context.Context) {
	ctx = context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(Config.Github.Token)},
	)
	client = oauth2.NewClient(ctx, ts)

	return
}
