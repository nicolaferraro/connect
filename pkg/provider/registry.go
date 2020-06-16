package provider

import (
	"github.com/nicolaferraro/connect/resources"
	"path"
)
import "gopkg.in/yaml.v2"

var providers []Provider

func init() {
	for _, id := range resources.Resources("/providers") {
		provDesc := resources.Resource(path.Join("/providers", id))
		var provider Provider
		if err := yaml.Unmarshal(provDesc, &provider); err != nil {
			panic(err)
		}
		for _, com := range CommonFields {
			if !provider.HasField(com.ID) {
				commProvField := com
				provider.Fields = append(provider.Fields, commProvField)
			}
		}
		providers = append(providers, provider)
	}
}

func Get(id string) *Provider {
	for _, p := range providers {
		if p.ID == id {
			provider := p
			return &provider
		}
	}
	return nil
}
