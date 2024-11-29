build:
	go build -o ./app ./cmd/main.go

test:
	go test -coverprofile=coverage.out ./...

docker-build:
	docker build -t usdt-app .

run:
	docker compose up -d

lint:
	golangci-lint run

docker-compose-up:
	docker compose up -d

docker-compose-down:
	docker compose down