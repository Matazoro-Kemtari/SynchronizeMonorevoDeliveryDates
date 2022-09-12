go build -ldflags "-X main.version=v1.1.5 -X main.revision=$(git rev-parse --short HEAD)" .
