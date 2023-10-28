FROM golang:1.20 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o gitecho ./cmd

# Use a minimal Alpine-based image as the production image
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/gitecho .

# Set the entry point for the container
ENTRYPOINT ["./gitecho"]
CMD ["-f", "config.yaml"]
