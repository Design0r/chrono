FROM golang:alpine AS build

ENV CGO_ENABLED=1
WORKDIR /app

RUN apk add --no-cache gcc musl-dev ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/build/chrono -ldflags="-s -w" ./cmd/main.go


FROM alpine:latest

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=build /app/build/chrono .
ENTRYPOINT ["./chrono"]
