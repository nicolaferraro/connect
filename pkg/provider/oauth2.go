package provider

import (
	"golang.org/x/oauth2"
	"strings"
)

func (provider Provider) GetOauth2Configuration() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     provider.GetField(ClientIDField),
		ClientSecret: provider.GetField(ClientSecretField),
		Endpoint: oauth2.Endpoint{
			AuthURL:  provider.GetField(AuthURLField),
			TokenURL: provider.GetField(TokenURLField),
		},
		RedirectURL: provider.GetField(RedirectURLField),
		Scopes:      strings.Split(provider.GetField(ScopesField), " "),
	}
}
