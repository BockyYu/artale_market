FROM golang:1.21-alpine

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "/app/entrypoint.sh"]
