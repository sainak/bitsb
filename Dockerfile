# Builder
FROM golang:1.20-alpine as builder

RUN apk update && apk upgrade && \
    apk --update add git make bash build-base

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata curl

WORKDIR /app 

COPY --from=builder /app/engine /app/

HEALTHCHECK --interval=1m30s --timeout=10s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:${WEBSITE_PORT}/ping || exit 1

CMD /app/engine