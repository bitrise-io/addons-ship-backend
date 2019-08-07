FROM quay.io/bitriseio/bitrise-base
RUN go get github.com/codegangsta/gin \
    && go get github.com/kisielk/errcheck \
    && go get -u golang.org/x/lint/golint
ADD . /bitrise/src
ARG GOFLAGS
ENV GOFLAGS $GOFLAGS
