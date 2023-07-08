FROM golang:1.20 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/server

# Use a minimal Alpine-based image as the production image
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/app .

# Set the entry point for the container
ENTRYPOINT ["./app"]
EXPOSE 8080
