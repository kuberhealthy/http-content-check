IMAGE := "kuberhealthy/http-content-check"
TAG := "latest"

# Build the http content check container locally.
build:
	podman build -f Containerfile -t {{IMAGE}}:{{TAG}} .

# Run the unit tests for the http content check.
test:
	go test ./...

# Build the http content check binary locally.
binary:
	go build -o bin/http-content-check ./cmd/http-content-check
