FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.20
WORKDIR /app

RUN adduser -D -H appuser
COPY --from=builder /app/server /app/server

USER appuser
EXPOSE 8080
CMD ["/app/server"]