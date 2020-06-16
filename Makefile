default: build

build: generate
	go build ./cmd/connect/

generate:
	./hack/embed_resources.sh ./resources
