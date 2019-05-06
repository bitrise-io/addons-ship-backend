FROM golang:1.12
RUN go get github.com/codegangsta/gin
WORKDIR /src
COPY . /src
RUN go mod download