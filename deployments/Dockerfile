# Builder
FROM golang:1.20-alpine as builder

RUN apk update && apk upgrade && \
    apk --update add git make bash build-base

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION="nil"

RUN VERSION="${VERSION}" make build

# Distribution
FROM alpine:latest as app

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata curl

WORKDIR /app

COPY --from=builder /app/build/ /app/

HEALTHCHECK --interval=1m --timeout=10s --start-period=3m --retries=3 \
    CMD curl -f http://localhost:${WEBSITE_PORT}/ping || exit 1

CMD /app/engine