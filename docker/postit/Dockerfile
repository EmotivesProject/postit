FROM golang:1.16.3-stretch

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

EXPOSE 80

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build ./cmd/api/" -command="./api"