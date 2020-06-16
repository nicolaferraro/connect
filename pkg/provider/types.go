package provider

import "fmt"

const (
	ClientIDField     = "client-id"
	ClientSecretField = "client-secret"
	ScopesField       = "scopes"
	RedirectURLField  = "redirect-url"
	AuthURLField      = "auth-url"
	TokenURLField     = "token-url"
)

var CommonFields = []ProviderField{
	{
		ID:          ClientIDField,
		Description: "Client ID of the registered application",
		Required:    true,
		Global:      true,
	},
	{
		ID:          ClientSecretField,
		Description: "Client secret of the registered application",
		Required:    true,
		Global:      true,
	},
	{
		ID:          ScopesField,
		Description: "Application specific scopes required for the authorization request",
		Required:    true,
	},
	{
		ID:          RedirectURLField,
		Description: "Redirect URL as configured in the remote application",
		Value:       "http://localhost:3000/callback",
	},
	{
		ID:          AuthURLField,
		Description: "Oauth2 Auth URL",
		Required:    true,
	},
	{
		ID:          TokenURLField,
		Description: "Oauth2 Token URL",
		Required:    true,
	},
}

type Provider struct {
	ID     string          `json:"id,omitempty"`
	Group  string          `json:"group,omitempty"`
	Fields []ProviderField `json:"fields,omitempty"`
}

func (p *Provider) SetField(id, value string) error {
	for idx, f := range p.Fields {
		if f.ID == id {
			p.Fields[idx].Value = value
			return nil
		}
	}
	return fmt.Errorf("unknown field %q", id)
}

func (p *Provider) HasField(id string) bool {
	for _, f := range p.Fields {
		if f.ID == id {
			return true
		}
	}
	return false
}

func (p *Provider) GetField(id string) string {
	for _, f := range p.Fields {
		if f.ID == id {
			return f.Value
		}
	}
	return ""
}

type ProviderField struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Value       string `json:"value,omitempty"`
	Global      bool   `json:"global,omitempty"`
}

type Oauth2Config struct {
	AuthURL  string `json:"authURL,omitempty"`
	TokenURL string `json:"tokenURL,omitempty"`
}
