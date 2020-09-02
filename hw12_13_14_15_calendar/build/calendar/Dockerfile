FROM golang:latest AS build

RUN mkdir -p /tmp/calendar
WORKDIR /tmp/calendar

COPY go.sum .
COPY go.mod .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /tmp/calendar/calendar cmd/calendar/main.go

FROM alpine:latest

COPY --from=build /tmp/calendar/calendar /app/calendar

RUN mkdir /app/logs
WORKDIR /app

CMD ["./calendar", "--config=main.json"]