#!/bin/bash

ROOT=$(git rev-parse --show-toplevel)

# docker run -it -v $ROOT:/go/src/go-anonymize-mysqldump golang:1.13.15-stretch bash

# Test tool
docker run -e GOOS=linux -e GOARCH=amd64 -v $ROOT:/go/src/go-anonymize-mysqldump golang:1.13.15-stretch bash -c \
	"apt update && apt install ca-certificates libgnutls30 -y \
	&& cd /go/src/go-anonymize-mysqldump \
	&& go get -u golang.org/x/sys || true && git -C \$(go env GOPATH)/src/golang.org/x/sys reset --hard 27713097b956 \
	&& go get -u github.com/stretchr/testify \
	&& go get . && go test"
