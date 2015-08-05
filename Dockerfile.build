FROM golang:1.4-cross
RUN go get github.com/tools/godep
ADD . /go/src/github.com/cpuguy83/docker-compose-watcher
WORKDIR /go/src/github.com/cpuguy83/docker-compose-watcher
ENV GOOS darwin
CMD ./make.sh
