FROM golang:1.24-alpine AS build

ENV CGO_ENABLED=1
RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev \
    make \
    nodejs \
    npm


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make docker-install
RUN make build


FROM scratch

WORKDIR /app
COPY --from=build /app/build/chrono .
COPY .env .

ENTRYPOINT ["./chrono"]
