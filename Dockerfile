FROM golang:1.24-bookworm AS build

RUN apt update && apt install -y make
RUN apt install nodejs -y
RUN apt install npm -y
RUN apt-get install -y ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make docker-install
RUN make build


FROM scratch

WORKDIR /app
COPY --from=build /app/build/Chrono .
COPY .env .

ENTRYPOINT ["./Chrono"]
