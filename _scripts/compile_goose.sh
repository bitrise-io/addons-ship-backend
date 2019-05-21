#!/bin/bash
set -ex

go build -i -o db/goose db/*.go
