go build -ldflags "-X main.version=v1.1.0 -X main.revision=$(git rev-parse --short HEAD)" .
