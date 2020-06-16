package authorization

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nicolaferraro/connect/pkg/provider"
	"github.com/nicolaferraro/connect/pkg/token"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

type Flow struct {
	provider.Provider
}

func NewFlow(provider provider.Provider) Flow {
	return Flow{provider}
}

func (f *Flow) RequestToken(ctx context.Context) (*token.Token, error) {
	state := uuid.New().String()

	server := NewServer(state)
	server.Start()

	oauth2Config := f.GetOauth2Configuration()

	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("A browser window is being opened on the following URL to proceed with authorization: %v\n", url)

	if err := browser.OpenURL(url); err != nil {
		fmt.Printf("ERROR: cannot open URL in browser: %v\n", err)
	}

	var code string
	select {
	case code = <-server.Code:
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout while waiting for code")
	}

	tk, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return &token.Token{
		Provider: f.Provider,
		Oauth2:   tk,
	}, nil
}
