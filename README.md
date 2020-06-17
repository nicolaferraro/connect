# Connect

Connect is a (badly named) POC that allows to authorize towards external Oauth 2.0 endpoints, 
store the credentials on a Kubernetes secret and keep it updated (Oauth 2.0 refresh token) via an agent.

It's currently a POC for my home automation routes, but can be extended with multiple providers (see `resources/providers`).

To create an initial secret:

- Create an application on a remote SaaS (e.g. Azure AD) and get client-id and client-secret (set callback URL to `http://localhost:3000`)
- Make sure you're logged in onto a Kubernetes cluster
- ./connect login azure -o client-id=xxx -o client-secret=yyy -o scopes="a b c"
- A browser window is opened to let you authorize
- When you authorize, a secret containing the tokens is stored on the current namespace on Kubernetes
- You can deploy the agent to the current namespace via `ko apply -f config/`
- The agent will refresh all your tokens whenever they expire   
