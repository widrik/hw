FROM golang:latest as builder

RUN mkdir -p /tmp/scheduler
WORKDIR /tmp/scheduler

COPY go.sum .
COPY go.mod .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /tmp/scheduler/scheduler cmd/scheduler/main.go

FROM alpine:latest

COPY --from=builder /tmp/scheduler/scheduler /app/scheduler

RUN mkdir /app/logs
WORKDIR /app

CMD ["./scheduler", "--config=scheduler.json"]