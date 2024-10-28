# Use a Go base image
FROM golang:1.23
EXPOSE 3000
WORKDIR /app
COPY . .
CMD ["go", "run", "cmd/server/main.go"]