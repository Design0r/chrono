FROM golang:1.23-bookworm AS build

RUN apt update && apt install -y make
RUN apt install nodejs -y
RUN apt install npm -y

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make docker-install
RUN make build


FROM ubuntu:latest

WORKDIR /app
COPY --from=build /app/build/Chrono .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY .env .

ENTRYPOINT ["./Chrono"]
