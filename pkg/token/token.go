package token

import (
	"context"
	"github.com/nicolaferraro/connect/pkg/provider"
	"golang.org/x/oauth2"
	"time"
)

type Token struct {
	Provider provider.Provider `json:"provider"`
	Oauth2   *oauth2.Token     `json:"oauth2"`
}

func (t *Token) Source(ctx context.Context) *TokenSource {
	oauthTokenSource := t.Provider.GetOauth2Configuration().TokenSource(ctx, t.Oauth2)
	return &TokenSource{Oauth2: oauthTokenSource}
}

func (t *Token) GetAccessToken() string {
	if t.Oauth2 != nil {
		return t.Oauth2.AccessToken
	}
	return ""
}

func (t *Token) GetExpiry() time.Time {
	if t.Oauth2 != nil {
		return t.Oauth2.Expiry
	}
	return time.Unix(1<<63-62135596801, 999999999)
}

func (t *Token) Refresh() (*Token, error) {
	source := t.Provider.GetOauth2Configuration().TokenSource(context.Background(), t.Oauth2)
	otk, err := source.Token()
	if err != nil {
		return nil, err
	}
	return &Token{
		Provider: t.Provider,
		Oauth2:   otk,
	}, nil
}

type TokenSource struct {
	Oauth2 oauth2.TokenSource
}
