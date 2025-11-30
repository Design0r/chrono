FROM golang:1.25-alpine AS build

ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/build/chrono -ldflags="-s -w" ./cmd/main.go

# Runtime image
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=build /app/build/chrono .
ENTRYPOINT ["./chrono"]
