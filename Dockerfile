FROM golang:alpine AS builder

WORKDIR /src
ADD . .
RUN go build -o sysl-catalog .

FROM alpine:latest
COPY --from=builder /src/sysl-catalog /bin/

ENTRYPOINT [ "sysl-catalog" ]
