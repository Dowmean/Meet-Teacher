FROM golang:alpine

RUN apk add build-base
RUN mkdir /app

WORKDIR /app
ADD Meet-Teacher/. .
RUN go mod download
RUN go install -mod=mod github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main