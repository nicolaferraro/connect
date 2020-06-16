default: build

build: generate build-connect build-agent

build-connect:
	go build ./cmd/connect/

build-agent:
	go build ./cmd/agent/

generate:
	./hack/embed_resources.sh ./resources
