FROM quay.io/bitriseio/bitrise-base:latest
RUN go get github.com/codegangsta/gin \
    && go get github.com/kisielk/errcheck \
    && go get -u golang.org/x/lint/golint
ADD . /bitrise/src
RUN go mod download
