FROM golang:1.27-alpine AS build

ENV CGO_ENABLED=1
WORKDIR /app

RUN apk add --no-cache gcc musl-dev ca-certificates

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY cmd ./cmd
COPY config ./config
COPY db ./db
COPY internal ./internal

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /app/build/chrono ./cmd/main.go


FROM alpine:3.24

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=build /app/build/chrono .
ENTRYPOINT ["./chrono"]
