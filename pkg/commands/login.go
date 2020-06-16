package commands

import (
	"context"
	"fmt"
	"github.com/nicolaferraro/connect/pkg/authorization"
	"github.com/nicolaferraro/connect/pkg/storage/kubernetes"
	kubernetesutils "github.com/nicolaferraro/connect/pkg/util/kubernetes"
	"github.com/spf13/cobra"
)

type loginOptions struct {
	Name      string
	Options   []string
	Namespace string
}

func NewCmdRegister() *cobra.Command {
	options := loginOptions{}

	cmd := cobra.Command{
		Use:   "login",
		Short: "Login to a cloud service and store credentials on Kubernetes",
		RunE:  options.login,
	}

	cmd.Flags().StringVar(&options.Name, "name", "", `The name to use for the token`)
	cmd.Flags().StringVar(&options.Namespace, "namespace", "", `The namespace to use for storing the token`)
	cmd.Flags().StringArrayVarP(&options.Options, "option", "o", nil, `A list of options to pass`)

	return &cmd
}

func (o *loginOptions) login(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("wrong syntax. Expected: %s %s <provider>", cmd.Parent().Name(), cmd.Name())
	}
	providerName := args[0]

	provider, err := loadProvider(providerName, o.Options)
	if err != nil {
		return err
	}

	flow := authorization.NewFlow(*provider)
	tk, err := flow.RequestToken(context.Background())

	namespace := o.Namespace
	if namespace == "" {
		namespace, err = kubernetesutils.GetCurrentNamespace()
		if err != nil {
			return err
		}
	}

	store, err := kubernetes.NewKubernetesTokenStorage(namespace)
	if err != nil {
		return err
	}

	name := o.Name
	if name == "" {
		name = providerName
	}

	if err = store.Save(name, tk); err != nil {
		return err
	}

	fmt.Printf("Token %q stored successfully\n", name)
	return nil
}
