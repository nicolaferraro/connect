package commands

import (
	"fmt"
	"github.com/nicolaferraro/connect/pkg/storage/kubernetes"
	kubernetesutils "github.com/nicolaferraro/connect/pkg/util/kubernetes"
	"github.com/spf13/cobra"
)

type refreshOptions struct {
	Namespace string
}

func NewCmdRefresh() *cobra.Command {
	options := refreshOptions{}

	cmd := cobra.Command{
		Use:   "refresh",
		Short: "Refresh a token already saved on Kubernetes",
		RunE:  options.refresh,
	}

	cmd.Flags().StringVar(&options.Namespace, "namespace", "", `The namespace to use when looking for the token`)

	return &cmd
}

func (o *refreshOptions) refresh(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("wrong syntax. Expected: %s %s <name>", cmd.Parent().Name(), cmd.Name())
	}
	name := args[0]

	namespace := o.Namespace
	if namespace == "" {
		ns, err := kubernetesutils.GetCurrentNamespace()
		if err != nil {
			return err
		}
		namespace = ns
	}

	store, err := kubernetes.NewKubernetesTokenStorage(namespace)
	if err != nil {
		return err
	}

	tk, err := store.Get(name)
	if err != nil {
		return err
	}

	newToken, err := tk.Refresh()
	if err != nil {
		return err
	}
	if newToken.GetAccessToken() == tk.GetAccessToken() {
		fmt.Printf("Token %q has not been refreshed. Refresh deadline is %v\n", name, newToken.GetExpiry())
		return nil
	} else {
		fmt.Printf("Token %q successfully refreshed\n", name)
		return nil
	}
}
