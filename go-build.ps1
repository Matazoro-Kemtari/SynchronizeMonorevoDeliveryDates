go build -ldflags "-X main.version=v1.0.2 -X main.revision=$(git rev-parse --short HEAD)" .
