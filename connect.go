package main

import (
	"context"
	"fmt"
	"github.com/nicolaferraro/connect/authorization"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
	"os"
	"strings"
	"time"
)

var Providers = map[string]authorization.Configuration{
	"azure": {
		Oauth2: &oauth2.Config{
			ClientID:     os.Getenv("AZURE_CLIENT_ID"),
			ClientSecret: os.Getenv("AZURE_CLIENT_SECRET"),
			Scopes:       strings.Split(os.Getenv("AZURE_SCOPES"), " "),
			Endpoint:     microsoft.AzureADEndpoint(os.Getenv("AZURE_TENANT")),
			RedirectURL:  "http://localhost:3000/callback",
		},
	},
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	authFlow := authorization.Flow{
		Providers["azure"],
	}

	token, err := authFlow.RequestToken(ctx)
	if err != nil {
		fmt.Printf("An error occurred: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Token: %v\n", token.Oauth2.AccessToken)

	tokenSource := token.Source(context.TODO())
	for i := 0; i < 2; i++ {
		newToken, err := tokenSource.Oauth2.Token()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Token (%s) expires %v\n", newToken.AccessToken[len(newToken.AccessToken)-5:], newToken.Expiry)
		time.Sleep(5 * time.Second)
	}

}
