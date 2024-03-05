FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config.env .
COPY --from=builder /app/migrate /usr/bin/migrate
COPY --from=builder /app/db/migration ./db/migration
COPY --from=builder /app/start.sh .
COPY --from=builder /app/wait-for.sh .

RUN chmod +x /app/start.sh

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]
