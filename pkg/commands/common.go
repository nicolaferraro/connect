package commands

import (
	"bytes"
	"fmt"
	"github.com/nicolaferraro/connect/pkg/provider"
	"os"
	"strings"
	"text/template"
)

func loadProvider(providerName string, options []string) (*provider.Provider, error) {
	provider := provider.Get(providerName)
	if provider == nil {
		return nil, fmt.Errorf("Provider %q not found", providerName)
	}

	for _, opt := range options {
		if !strings.Contains(opt, "=") {
			return nil, fmt.Errorf("wrong format for option %q: expected <key>=<val>", opt)
		}
		kvs := strings.SplitN(opt, "=", 2)
		if err := provider.SetField(kvs[0], kvs[1]); err != nil {
			return nil, err
		}
	}

	// Filling with env vars if missing
	for idx := range provider.Fields {
		if provider.Fields[idx].Value == "" {
			provider.Fields[idx].Value = getOptionAsEnv(provider.ID, provider.Fields[idx].ID)
		}
		if provider.Fields[idx].Value == "" {
			provider.Fields[idx].Value = getOptionAsEnv(provider.Group, provider.Fields[idx].ID)
		}
	}

	for _, o := range provider.Fields {
		if o.Required && o.Value == "" {
			return nil, fmt.Errorf("required field %q is missing", o.ID)
		}
	}

	params := make(map[string]string)
	for _, o := range provider.Fields {
		params[o.ID] = o.Value
	}

	for idx, o := range provider.Fields {
		t, err := template.New(o.ID).Parse(o.Value)
		if err != nil {
			return nil, err
		}
		res := bytes.NewBuffer(nil)
		err = t.Execute(res, params)
		if err != nil {
			return nil, err
		}
		provider.Fields[idx].Value = res.String()
	}

	return provider, nil
}

func getOptionAsEnv(key, option string) string {
	name := fmt.Sprintf("%s_%s", key, option)
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ToUpper(name)
	return os.Getenv(name)
}
