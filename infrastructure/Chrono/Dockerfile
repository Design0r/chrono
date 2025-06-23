FROM golang:1.24-alpine AS build

ENV CGO_ENABLED=1
RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev \
    nodejs \
    npm \
    ca-certificates


WORKDIR /app

# install go dependnecies
COPY go.mod go.sum ./
RUN go mod download
RUN	go install github.com/a-h/templ/cmd/templ@latest

# install nodejs dependnecies
COPY package.json .
RUN npm install

COPY . .

RUN npm run build
RUN	templ generate
RUN	go build -o ./build/chrono -ldflags='-s -w -extldflags "-static"' ./cmd/main.go

FROM scratch

WORKDIR /app
COPY --from=build /app/build/chrono .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY .env .

ENTRYPOINT ["./chrono"]
