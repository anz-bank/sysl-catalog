# Build sysl-catalog binary
FROM golang:alpine AS builder
ADD . /src
RUN cd /src && CGO_ENABLED=0 go build -o sysl-catalog

FROM alpine:latest
RUN apk add --no-cache chromium
WORKDIR /usr
COPY --from=builder /src/sysl-catalog /bin/

ENTRYPOINT [ "sysl-catalog" ]
