FROM quay.io/bitriseio/bitrise-base
RUN go get github.com/codegangsta/gin
ADD . /bitrise/src
ARG GITHUB_TOKEN
RUN echo $GITHUB_TOKEN
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com".insteadOf "https://github.com"
RUN go mod download