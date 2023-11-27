E2E_TEST_PKGs := ./pkg/auth
E2E_TEST_OUTPUT := e2e_coverage.out

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
test-cover:
	find . -type f -name '*_test.go' -exec dirname {} \; | sort -u | xargs go test -cover

# Install linter dependencies
lint-dep:
	go install github.com/daixiang0/gci@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/bombsimon/wsl/v4/cmd...@master

# Start the DB container
db-up:
	docker compose up -d postgres-db

# Stop and remove the DB container
db-down:
	docker compose stop postgres-db
	docker compose rm -f postgres-db

# Run tests
test:
	find . -type f -name '*_test.go' -not -path './e2e/*' -exec dirname {} \; | sort -u | xargs go test -cover
	go test -cover -coverpkg=$(E2E_TEST_PKGs) ./e2e

.PHONY: e2e
e2e:
	go test -cover -coverpkg=$(E2E_TEST_PKGs) -coverprofile=$(E2E_TEST_OUTPUT) ./e2e
	go tool cover -html=$(E2E_TEST_OUTPUT)