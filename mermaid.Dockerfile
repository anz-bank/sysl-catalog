# Build sysl-catalog binary
FROM golang:alpine AS builder

# Caches go dependencies in the first layer 
# Speeds up builds when go.mod and go.sum is unchanged
RUN mkdir /src
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download

ADD . /src
RUN cd /src && CGO_ENABLED=0 go build -o sysl-catalog

FROM alpine:latest
RUN apk add --no-cache chromium
WORKDIR /usr
COPY --from=builder /src/sysl-catalog /bin/

ENTRYPOINT [ "sysl-catalog" ]
