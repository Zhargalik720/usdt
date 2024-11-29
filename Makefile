build:
	go build -o ./app ./cmd/main.go

test:
	go test -coverprofile=coverage.out ./...

docker-build:
	docker build -t usdt-app .

run: build
	./docker_compose_manager.sh

lint:
	golangci-lint run

docker-compose-up:
	./docker_compose_manager.sh

docker-compose-down:
	docker compose down