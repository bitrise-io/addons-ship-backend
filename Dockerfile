FROM quay.io/bitriseio/bitrise-base
RUN rm -rf /usr/local/go
RUN wget -q https://storage.googleapis.com/golang/go1.12.linux-amd64.tar.gz -O go-bins.tar.gz && tar -C /usr/local -xvzf go-bins.tar.gz && rm go-bins.tar.gz
RUN go get github.com/codegangsta/gin
RUN go mod download