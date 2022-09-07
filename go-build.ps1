go build -ldflags "-X main.version=v1.0.1 -X main.revision=$(git rev-parse --short HEAD)" .
