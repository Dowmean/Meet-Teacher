# Build stage
FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Run stage
FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]
