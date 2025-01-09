FROM golang:1.23-bookworm AS build

RUN apt update && apt install -y make
RUN apt install nodejs -y
RUN apt install npm -y

WORKDIR /app

COPY . .

RUN make docker-install
RUN make build


FROM ubuntu:latest

WORKDIR /app
COPY --from=build /app/build/Chrono .
COPY .env .

ENTRYPOINT ["./Chrono"]
