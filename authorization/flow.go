package authorization

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type Flow struct {
	Configuration Configuration
}

func (f *Flow) RequestToken(ctx context.Context) (*Token, error) {
	state := uuid.New().String()

	server := NewServer(state)
	server.Start()

	url := f.Configuration.Oauth2.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	var code string
	select {
	case code = <-server.Code:
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout while waiting for code")
	}

	token, err := f.Configuration.Oauth2.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return &Token{
		Configuration: f.Configuration,
		Oauth2:        token,
	}, nil
}

type Configuration struct {
	Oauth2 *oauth2.Config
}

type Token struct {
	Configuration Configuration `json:"configuration"`
	Oauth2        *oauth2.Token `json:"oauth2"`
}

func (t *Token) Source(ctx context.Context) *TokenSource {
	oauthTokenSource := t.Configuration.Oauth2.TokenSource(ctx, t.Oauth2)
	return &TokenSource{Oauth2: oauthTokenSource}
}

type TokenSource struct {
	Oauth2 oauth2.TokenSource
}
