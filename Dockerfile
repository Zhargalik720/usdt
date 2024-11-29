FROM golang:1.23-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod tidy
COPY . .

RUN go build -o /app/main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/.env /.env
COPY --from=builder /app/main /app/main

COPY --from=builder /app/internal/infrastructure/db/migrations /app/migrations

CMD ["/app/main"]
