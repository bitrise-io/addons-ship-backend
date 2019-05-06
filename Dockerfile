FROM golang:1.12
WORKDIR /src
COPY . /src
RUN go mod download