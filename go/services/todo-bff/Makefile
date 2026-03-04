SERVICE_NAME=todo-bff
GO=go
SERVICE_PORT=5006

.PHONY: build
build:
	go mod download
	go build -o dist/main ./cmd/$(SERVICE_NAME)/...

run:
	export $(shell cat .local.env | xargs) && go run cmd/todo-bff/main.go

.PHONY: gqlgen
gqlgen:
	$(GO) run github.com/99designs/gqlgen generate

.PHONY: mockgen
mockgen:
	mockgen -destination=internal/domain/gateway/mock/gateway.go -source=internal/domain/gateway/gateway.go
	mockgen -destination=internal/usecase/mock/usecase.go -source=internal/usecase/usecase.go

.PHONY: wiregen
wiregen:
	go get github.com/google/wire/cmd/wire
	cd internal/registry; wire
	go mod tidy

generate: gqlgen mockgen wiregen

test-local:
	go install github.com/tenntenn/testtime/cmd/testtime@latest
	export $(shell cat local.test.env | xargs) && go test -timeout 30m -overlay=`testtime` ./...

test-local-integration:
	go install github.com/tenntenn/testtime/cmd/testtime@latest
	export $(shell cat .local.test.env | xargs) && go test -timeout 30m -overlay=`testtime` ./test/integration

lint:
	golangci-lint run

pre-push: lint test-local ## Run lint and test before push
