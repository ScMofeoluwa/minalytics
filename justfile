# Load environment variables from a .env file
set dotenv-load

# Create a new migration file
# Usage: just create-migration <migration_name>
create-migration m:
  migrate create -ext sql -dir database/migrations -seq {{m}}

# Apply database migrations
migrate-up:
  migrate -path database/migrations -database $DATABASE_URL -verbose up

# Rollback database migrations
migrate-down:
  migrate -path database/migrations -database $DATABASE_URL -verbose down

migrate-fix v:
  migrate -path database/migrations -database $DATABASE_URL force {{v}}

sqlc:
  sqlc generate

swag:
  swag init --generalInfo server/main.go --parseDependency --output ./docs

mock:
  mockery --name=AnalyticsService --filename=service.go --outpkg=mocks --with-expecter=true --recursive=true

test:
  go test -v ./...

coverage:
  go test -v $(go list ./... | grep -v -E '/(cmd|mocks|docs)') -covermode=atomic -coverprofile=coverage.out && go tool cover -html=coverage.out
  # go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
