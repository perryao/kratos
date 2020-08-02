package oidc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ory/herodot"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type ProviderGitLab struct {
	*ProviderGenericOIDC
}

func NewProviderGitLab(
	config *Configuration,
	public *url.URL,
) *ProviderGitLab {
	return &ProviderGitLab{
		ProviderGenericOIDC: &ProviderGenericOIDC{
			config: config,
			public: public,
		},
	}
}

func (g *ProviderGitLab) Claims(ctx context.Context, exchange *oauth2.Token) (*Claims, error) {
	o, err := g.OAuth2(ctx)
	if err != nil {
		return nil, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("%s", err))
	}

	client := o.Client(ctx, exchange)
	req, err := http.NewRequest("GET", "https://gitlab.com/oauth/userinfo", nil)
	if err != nil {
		return nil, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("%s", err))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("%s", err))
	}
	defer resp.Body.Close()

	var claims Claims
	if err := json.NewDecoder(resp.Body).Decode(&claims); err != nil {
		return nil, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("%s", err))
	}

	return &claims, nil
}
