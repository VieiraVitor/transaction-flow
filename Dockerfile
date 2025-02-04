# Build
FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o transaction-flow ./cmd/transaction-flow

# Execute
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=build /app/transaction-flow /app/transaction-flow

RUN chmod +x /app/transaction-flow

EXPOSE 8080

ENTRYPOINT ["/app/transaction-flow"]