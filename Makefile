.PHONY: bindata fmt test

# Build the binary for the current platform
make:
	go build -o bin/deltaiota ./cmd/deltaiota/

# Build binary assets
bindata:
	go-bindata -pkg bindata -ignore deltaiota.sql -o ./bindata/bindata.go res/...

# Format, vet, and lint all files
fmt:
	go fmt ./...
	go vet ./...
	golint .

# Run all tests
test:
	go test -v ./...
