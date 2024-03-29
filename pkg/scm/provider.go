package scm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/buyoio/runr/pkg/scm/github"
)

type Provider interface {
	RunnerToken() (string, error)
	RunnerSetup() ([]string, error)
	RunnerVersion() (string, error)
	Platform() string
	Organization() string
	Repository() string
	Token() string
	TeamSSHKeys(string) ([]string, error)
	SetLabels([]string)

	RemoveRunner(string) error
}

type SCMPlatform struct {
	Platform     string
	Organization string
	Repository   *string
	Token        string
}

type ProviderOptions struct {
	SCMPlatform
	Labels []string
}

func NewProvider(ctx context.Context, logger *slog.Logger, o *ProviderOptions) (Provider, error) {
	var err error
	var provider Provider

	switch o.Platform {
	case "github":
		// todo check another way for the token
		provider, err = github.NewProvider(ctx, logger, github.Options{
			Owner:      o.Organization,
			Repository: o.Repository,
			Token:      o.Token,
			Labels:     o.Labels,
		})
		if err != nil {
			logger.ErrorContext(ctx, "failed to create github provider", "error", err)
			return nil, err
		}
		logger.InfoContext(ctx, "github provider created")
	default:
		logger.ErrorContext(ctx, "unsupported SCM platform",
			slog.String("platform", o.Platform),
		)
		return nil, fmt.Errorf("unsupported SCM platform: %s", o.Platform)
	}

	return provider, nil
}
