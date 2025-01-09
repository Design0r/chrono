FROM golang:1.23-bookworm AS build

RUN apt update && apt install -y make
RUN apt install nodejs -y
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

WORKDIR /app

COPY . .

CMD ["make", "install"]
CMD ["make", "build"]
CMD ["sleep", "5m"]


# FROM scratch
#
# WORKDIR /app
# COPY --from=build /app/build/Chrono.exe ./chrono
# COPY .env .
#
# ENTRYPOINT ["./chrono"]
