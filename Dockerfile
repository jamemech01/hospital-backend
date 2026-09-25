FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM alpine:3.22

RUN adduser -D -H appuser

WORKDIR /app

COPY --from=builder /app/app /app/app

USER appuser

EXPOSE 8080

CMD ["/app/app"]