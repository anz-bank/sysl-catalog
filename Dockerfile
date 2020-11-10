FROM golang:alpine AS builder

WORKDIR /src
COPY go.* /src/
RUN go mod download
ADD . .
RUN go build -o sysl-catalog .

FROM alpine:latest
WORKDIR /usr
COPY --from=builder /src/sysl-catalog /bin/

ENTRYPOINT [ "sysl-catalog" ]
