FROM golang:1.22.3

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download && go mod verify

COPY . /app/

RUN go build -v -o /usr/local/bin/hyperbloom -ldflags="-s -w" ./cmd/app

RUN go clean

ENTRYPOINT [ "hyperbloom" ]