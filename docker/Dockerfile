FROM golang:1.15.2-alpine AS builder

RUN apk add git gcc libc-dev && git clone https://github.com/mgr9525/gokins.git /build
WORKDIR /build
RUN GOOS=linux GOARCH=amd64 go build -o bin/gokins main.go


FROM alpine:latest AS final

RUN apk update \
    && apk upgrade \
    && apk --no-cache add openssl \
    && apk --no-cache add ca-certificates \
    && rm -rf /var/cache/apk \
    && mkdir -p /app

COPY --from=builder /build/bin/gokins /app
WORKDIR /app
ENTRYPOINT ["/app/gokins"]