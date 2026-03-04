SERVICE_NAME=auth
DOCKER_COMPOSE=docker compose
GO=go
SERVICE_PORT=5007

.PHONY: build
build:
	go mod download
	go build -o dist/main ./cmd/$(SERVICE_NAME)/...

run:
	@echo "Running $(SERVICE_NAME) service..."
	$(GO) run ./cmd/$(SERVICE_NAME)

.PHONY: mockgen
mockgen:
	mockgen -destination=internal/domain/gateway/mock/gateway.go -source=internal/domain/gateway/gateway.go
	mockgen -destination=internal/usecase/mock/usecase.go -source=internal/usecase/usecase.go

.PHONY: wiregen
wiregen:
	go get github.com/google/wire/cmd/wire
	cd internal/registry; wire
	go mod tidy

generate: mockgen wiregen

test-local:
	go install github.com/tenntenn/testtime/cmd/testtime@latest
	export $(shell cat .local.test.env | xargs) && go test -timeout 30m -overlay=`testtime` ./...

test-local-integration:
	go install github.com/tenntenn/testtime/cmd/testtime@latest
	export $(shell cat .local.test.env | xargs) && go test -timeout 30m -overlay=`testtime` ./test/integration

lint:
	golangci-lint run

pre-push: lint test-local ## Run lint and test before push

run:
	export $(shell cat .local.env | xargs) && go run cmd/auth/main.go

docker-up:
	@echo "Starting Docker containers for $(SERVICE_NAME)..."
	$(DOCKER_COMPOSE) up -d

docker-down:
	@echo "Stopping Docker containers for $(SERVICE_NAME)..."
	$(DOCKER_COMPOSE) down

docker-logs:
	@echo "Showing Docker logs for $(SERVICE_NAME)..."
	$(DOCKER_COMPOSE) logs -f

migrate-up:
	@echo "Running migrations up for $(SERVICE_NAME)..."
	$(GO) run ./cmd/migrate up

migrate-down:
	@echo "Running migrations down for $(SERVICE_NAME)..."
	$(GO) run ./cmd/migrate down

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration $(NAME)..."
	migrate create -ext sql -dir ./database/migrations -seq $(NAME)
