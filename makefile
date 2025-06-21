DB_URL=postgres://postgres:@localhost:5432/test_db_drone?sslmode=disable

run:
	go run cmd/main.go

swag-init:
	swag init -g internal/server/router.go --parseDependency true --parseInternal true --parseDepth 2 --output ./docs

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

migrate-new:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq "$$name"

mockgen:
	go generate ./...

test:
	go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out

test-coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out