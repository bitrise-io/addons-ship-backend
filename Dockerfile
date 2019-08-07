FROM quay.io/bitriseio/bitrise-base
RUN go get github.com/codegangsta/gin
ADD . /bitrise/src
ARG GOFLAGS
ENV GOFLAGS $GOFLAGS
