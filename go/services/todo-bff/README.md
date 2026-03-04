# Todo-BFF Service

## Getting Started

### 1. Setup Environment

Copy the environment file and configure it:

```bash
cp .local.env .env
```

Edit `.local.env` with your configuration:

```
ENV="local"
PORT="5006"
ALLOW_ORIGINS="*"
TODO_ADDR="localhost:5005"
AUTH_ADDR="localhost:5007"
```

### 2. Install Dependencies and Generate Code

```bash
make build
make generate
```

### 3. Run the Service

```bash
make run
```

Note: The `run` command automatically loads environment variables from `.local.env`

The GraphQL API will be available at `http://localhost:5006/graphql`

### 4. Access GraphQL Playground

Navigate to `http://localhost:5006/playground` in your browser to access the GraphQL Playground interface.

## Testing

### Run All Tests

```bash
make test-local
```

### Run Integration Tests

```bash
make test-local-integration
```

### Pre-Push Checks

Run linting and tests before pushing:

```bash
make pre-push
```

## Development

### Schema Changes

1. Edit the GraphQL schema in `internal/handler/graph/*.graphqls`
2. Regenerate code: `make generate`
3. Implement resolver logic in `internal/handler/graph/*.resolvers.go`

### Code Generation

#### Generate All Code

```bash
make generate
```

This generates GraphQL code, mocks, and wire dependencies.

#### Generate GraphQL Code Only

```bash
make gqlgen
```

#### Generate Mocks Only

```bash
make mockgen
```

#### Generate Wire Dependencies Only

```bash
make wiregen
```

### Code Quality

#### Run Linter

```bash
make lint
```

## GraphQL Playground

The service includes GraphQL Playground for interactive API exploration:

1. Start the service: `make run`
2. Navigate to: `http://localhost:5006/playground`
3. Explore the schema and test queries/mutations
