package commands

import (
	"context"
	"fmt"
	"github.com/nicolaferraro/connect/pkg/storage/kubernetes"
	"github.com/spf13/cobra"
)

type refreshOptions struct {
}

func NewCmdRefresh() *cobra.Command {
	options := refreshOptions{}

	cmd := cobra.Command{
		Use:   "refresh",
		Short: "Refresh a token already saved on Kubernetes",
		RunE:  options.refresh,
	}

	return &cmd
}

func (o *refreshOptions) refresh(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("wrong syntax. Expected: %s %s <name>", cmd.Parent().Name(), cmd.Name())
	}
	name := args[0]

	store, err := kubernetes.NewKubernetesTokenStorage("nicola-webhooks")
	if err != nil {
		return err
	}

	tk, err := store.Get(name)
	if err != nil {
		return err
	}

	source := tk.Provider.GetOauth2Configuration().TokenSource(context.Background(), tk.Oauth2)
	newToken, err := source.Token()
	if err != nil {
		return err
	}
	if newToken.AccessToken == tk.GetAccessToken() {
		fmt.Printf("Token %q has not been refreshed. Refresh deadline is %v\n", name, newToken.Expiry)
		return nil
	} else {
		fmt.Printf("Token %q successfully refreshed\n", name)
		return nil
	}
}
