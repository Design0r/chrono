FROM golang:1.25-alpine AS build

ENV CGO_ENABLED=1
RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev \
    ca-certificates


WORKDIR /app

# install go dependnecies
COPY go.mod go.sum ./

RUN go mod download


COPY . .

RUN	go build -o ./build/chrono -ldflags='-s -w -extldflags "-static"' ./cmd/main.go

FROM scratch

WORKDIR /app
COPY --from=build /app/build/chrono .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["./chrono"]
