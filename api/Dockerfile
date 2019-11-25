FROM golang:1.12.10 as build-env

LABEL maintainer="Damir Asylkhan<krosava_15d@mail.ru>"

WORKDIR /api

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /opt/api

# Release
FROM alpine:latest

WORKDIR /root/

COPY --from=build-env /opt/api .
COPY --from=build-env /api/config ./config

CMD ["./api", "run", "--config=./config/development/config.toml"]