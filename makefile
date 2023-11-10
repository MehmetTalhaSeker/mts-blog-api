# Run the linter on the specified path
lint:
	go mod tidy
	go vet ./...
	go fmt ./...
	gci write -s standard -s default -s "prefix(github.com/MehmetTalhaSeker/mts-blog-api)" .
	gofumpt -l -w .
	wsl -fix ./...
	golangci-lint run $(p)
.PHONY: lint

# Run and cover every tests
test:
	find . -type f -name '*_test.go' -exec dirname {} \; | sort -u | xargs go test -cover

# Install linter dependencies
lint-dep:
	go install github.com/daixiang0/gci@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/bombsimon/wsl/v4/cmd...@master


