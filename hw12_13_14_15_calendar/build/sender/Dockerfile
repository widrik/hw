FROM golang:latest as builder

RUN mkdir -p /tmp/sender
WORKDIR /tmp/sender

COPY go.sum .
COPY go.mod .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /tmp/sender/sender cmd/sender/main.go

FROM alpine:latest

COPY --from=builder /tmp/sender/sender /app/sender

RUN mkdir /app/logs
WORKDIR /app

CMD ["./sender", "--config=sender.json"]