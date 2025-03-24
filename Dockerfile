# Go Version
FROM golang:1.23

ENV GO111MODULE=on \
    GOFLAGS=-buildvcs=false

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.51.0

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]