.PHONY: dev test lint migrate-up migrate-down seed run docker-up docker-down

dev:
	go run ./cmd/server

test:
	go test -v -race -cover ./...

test-integration:
	go test -v -tags=integration ./...

lint:
	go vet ./...
	go fmt ./...

migrate-up:
	migrate -path migrations -database "sqlite3://./roombooker.db" up

migrate-down:
	migrate -path migrations -database "sqlite3://./roombooker.db" down

seed:
	migrate -path migrations -database "sqlite3://./roombooker.db" up 2

run:
	go build -o bin/server ./cmd/server
	./bin/server

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

clean:
	rm -f bin/server roombooker.db
