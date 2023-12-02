# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copying the source code
COPY *.go ./

# Build
RUN go build -o /squash-slack-bot

EXPOSE 8080

CMD ["/squash-slack-bot"]