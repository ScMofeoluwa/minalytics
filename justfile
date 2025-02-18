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

mock:
  mockery --name=Querier --filename=querier.go --outpkg=mocks --with-expecter=true --recursive=true

test:
  go test -v ./...
