FROM golang:1.12.10

LABEL maintainer="Damir Asylkhan<krosava_15d@mail.ru>"

WORKDIR /opt/e2e

COPY go.mod go.sum ./

RUN go mod download

COPY . .

CMD ["go", "test", "-v", "./..."]