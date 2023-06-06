FROM golang:1.20.4

COPY ./go.mod ./go.sum *.go ./src/
COPY ./example ./src/example

WORKDIR src

RUN go build -o /bin/example-app example/main.go

ENTRYPOINT /bin/example-app
