package token

import (
	"context"
	"github.com/nicolaferraro/connect/pkg/provider"
	"golang.org/x/oauth2"
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

type TokenSource struct {
	Oauth2 oauth2.TokenSource
}
