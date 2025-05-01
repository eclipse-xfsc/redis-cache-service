FROM golang:1.24.2-alpine3.20 as builder

RUN apk add git

WORKDIR /go/src/github.com/eclipse-xfsc/redis-cache-service

ARG APP_REPO_TAG

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.Version=$APP_REPO_TAG" -mod=vendor -o /tmp/cache ./cmd/cache/...

FROM alpine:3.19 as runner

COPY --from=builder /tmp/cache /opt/cache

WORKDIR /opt

CMD ["./cache"]
