# Todo Service

## Getting Started

### 1. Setup Environment

Copy the environment file and configure it:

```bash
cp .local.env .env
```

Edit `.local.env` with your configuration:

```
ENV="local"
SERVER_PORT=5005
DB_HOST=127.0.0.1
DB_PORT=33062
DB_USER=root
DB_PASS=password
DB_NAME=todo
GRPC_REFLECTION_ENABLE=true
```

### 2. Start the Database

```bash
make docker-up
```

This will start the MySQL database container on port 33062.

### 3. Run Migrations

```bash
make migrate-up
```

### 4. Run the Service

```bash
make run
```

Note: The `run` command automatically loads environment variables from `.local.env`

The service will start on port 5005 (or the port specified in your `.env` file).

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

## Database Migrations

### Create a New Migration

```bash
make migrate-create NAME=your_migration_name
```

This will create two files in `database/migrations/`:
- `XXXXXX_your_migration_name.up.sql` - Migration to apply
- `XXXXXX_your_migration_name.down.sql` - Migration to rollback

### Apply Migrations

```bash
make migrate-up
```

### Rollback Migrations

```bash
make migrate-down
```

## Development

### Code Generation

#### Generate All Code

```bash
make generate
```

This generates both mocks and wire dependencies.

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
