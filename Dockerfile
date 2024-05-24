FROM golang:1.22.2-alpine

WORKDIR /app

COPY ./cmd ./cmd
# COPY ./internal ./internal
# COPY ./pkg ./pkg
# COPY ./configs ./configs
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod tidy

RUN go build -o /node cmd/main.go

CMD /node