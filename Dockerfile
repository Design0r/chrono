FROM golang:alpine AS build

ENV CGO_ENABLED=1
WORKDIR /app

RUN apk add --no-cache gcc musl-dev ca-certificates

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY cmd ./cmd
COPY config ./config
COPY db/*.go ./db/
COPY db/migrations ./db/migrations
COPY db/repo ./db/repo
COPY internal ./internal

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /out/chrono -ldflags="-s -w" ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=build /out/chrono ./chrono
ENTRYPOINT ["./chrono"]
