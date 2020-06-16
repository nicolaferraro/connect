package commands

import (
	"bytes"
	"context"
	"fmt"
	"github.com/nicolaferraro/connect/pkg/authorization"
	"github.com/nicolaferraro/connect/pkg/provider"
	"github.com/nicolaferraro/connect/pkg/storage/kubernetes"
	"github.com/spf13/cobra"
	"strings"
	"text/template"
)

type loginOptions struct {
	Name    string
	Options []string
}

func NewCmdRegister() *cobra.Command {
	options := loginOptions{}

	cmd := cobra.Command{
		Use:   "login",
		Short: "Login to a cloud service and store credentials on Kubernetes",
		RunE:  options.login,
	}

	cmd.Flags().StringVar(&options.Name, "name", "", `The name to use for the token`)
	cmd.Flags().StringArrayVarP(&options.Options, "option", "o", nil, `A list of options to pass`)

	return &cmd
}

func (o *loginOptions) login(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("wrong syntax. Expected: %s %s <provider>", cmd.Parent().Name(), cmd.Name())
	}
	providerName := args[0]
	provider := provider.Get(providerName)
	if provider == nil {
		return fmt.Errorf("Provider %q not found", providerName)
	}

	for _, opt := range o.Options {
		if !strings.Contains(opt, "=") {
			return fmt.Errorf("wrong format for option %q: expected <key>=<val>", opt)
		}
		kvs := strings.SplitN(opt, "=", 2)
		if err := provider.SetField(kvs[0], kvs[1]); err != nil {
			return err
		}
	}

	for _, o := range provider.Fields {
		if o.Required && o.Value == "" {
			return fmt.Errorf("required field %q is missing", o.ID)
		}
	}

	params := make(map[string]string)
	for _, o := range provider.Fields {
		params[o.ID] = o.Value
	}

	for idx, o := range provider.Fields {
		t, err := template.New(o.ID).Parse(o.Value)
		if err != nil {
			return err
		}
		res := bytes.NewBuffer(nil)
		err = t.Execute(res, params)
		if err != nil {
			return err
		}
		provider.Fields[idx].Value = res.String()
	}

	flow := authorization.NewFlow(*provider)
	tk, err := flow.RequestToken(context.Background())

	store, err := kubernetes.NewKubernetesTokenStorage("nicola-webhooks")
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
