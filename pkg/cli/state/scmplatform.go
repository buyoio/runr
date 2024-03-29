package state

import (
	"context"
	"log/slog"

	"github.com/buyoio/runr/pkg/scm"
)

func (g *Github) SCMPlatform() *SCMPlatform {
	return &SCMPlatform{
		Platform:     "github",
		Organization: g.Organization,
		Repository:   g.Repository,
		Token:        g.Token,
	}
}

func (g *Gitlab) SCMPlatform() *SCMPlatform {
	return &SCMPlatform{
		Platform:     "gitlab",
		Organization: g.Organization,
		Repository:   g.Repository,
		Token:        g.Token,
	}
}

func (s *SCMPlatform) NewProvider(ctx context.Context, logger *slog.Logger) (scm.Provider, error) {
	provider, err := scm.NewProvider(ctx, logger, &scm.ProviderOptions{
		SCMPlatform: scm.SCMPlatform(*s),
	})
	if err != nil {
		return nil, err
	}
	return provider, nil
}
