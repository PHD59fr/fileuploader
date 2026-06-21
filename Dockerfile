FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY . .

RUN apk upgrade --no-cache && \
    apk add --no-cache upx && \
    go mod tidy && \
    CGO_ENABLED=0 go build -o torrentuploader . && \
    upx --best --lzma torrentuploader

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/torrentuploader .
COPY config.yaml .
COPY templates/ ./templates/

EXPOSE 8080

CMD ["./torrentuploader"]
