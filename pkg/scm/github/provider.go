package github

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/go-github/v59/github"
)

type Options struct {
	Owner      string
	Repository *string
	Token      string
	Labels     []string
}

type provider struct {
	ctx     context.Context
	logger  *slog.Logger
	options Options
	github  *github.Client

	runnerToken *github.RegistrationToken
}

func NewProvider(ctx context.Context, logger *slog.Logger, options Options) (*provider, error) {
	if options.Token == "" {
		return nil, fmt.Errorf("github token is required")
	}
	return &provider{
		ctx:     ctx,
		logger:  logger,
		options: options,
		github:  github.NewClient(nil).WithAuthToken(options.Token),
	}, nil
}

func (g *provider) Platform() string {
	return "github"
}

func (g *provider) Repository() string {
	if g.options.Repository == nil {
		return ""
	}
	return *g.options.Repository
}

func (g *provider) Organization() string {
	return g.options.Owner
}

func (g *provider) Token() string {
	return g.options.Token
}

func (g *provider) SetLabels(labels []string) {
	g.options.Labels = labels
}
