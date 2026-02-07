# syntax=docker/dockerfile:1

FROM golang:1.23.3-alpine AS builder

ARG CI_USER
ARG CI_TOKEN
ARG APP_VERSION="undefined"
ARG BUILD_TIME="undefined"

WORKDIR /go/src/github.com/artarts36/ip-scanner

RUN echo -e "machine github.com login $CI_USER password $CI_TOKEN" > ~/.netrc

RUN apk add git

COPY go.mod go.sum ./
COPY pkg/ip-scanner-grpc-api/go.mod pkg/ip-scanner-grpc-api/go.sum ./pkg/ip-scanner-grpc-api/
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X 'main.Version=${APP_VERSION}' -X 'main.BuildDate=${BUILD_TIME}'" -o /go/bin/ip-scanner /go/src/github.com/artarts36/ip-scanner/cmd/grpc/main.go

######################################################

FROM alpine

COPY --from=builder /go/bin/ip-scanner /go/bin/ip-scanner

WORKDIR app
COPY dbip.mmdb dbip.mmdb

# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.title="ip-scanner"
LABEL org.opencontainers.image.description="ip-scanner"
LABEL org.opencontainers.image.url="https://github.com/artarts36/ip-scanner"
LABEL org.opencontainers.image.source="https://github.com/artarts36/ip-scanner"
LABEL org.opencontainers.image.vendor="ArtARTs36"
LABEL org.opencontainers.image.version="$APP_VERSION"
LABEL org.opencontainers.image.created="$BUILD_TIME"
LABEL org.opencontainers.image.licenses="MIT"

EXPOSE 8000

CMD ["/go/bin/ip-scanner"]
