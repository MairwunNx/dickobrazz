FROM alpine:3.21 AS certs
RUN apk --no-cache add ca-certificates tzdata

FROM golang:1.25-alpine AS builder

RUN apk --no-cache add git
COPY go.mod go.sum ./
RUN go mod download

COPY program.go ./
COPY .version ./
COPY .git ./
COPY application/ ./application/

RUN mkdir -p /app
RUN go build -o /app/dickobot \
    -a -installsuffix cgo \
    -gcflags "all=-N -l" \
    -tags timetzdata \
    -ldflags="-s -w \
        -X dickobot/application/logging.Version=$(cat .version) \
        -X dickobot/application/logging.BuildAt=$(date +%Y-%m-%d_%H:%M:%S) \
        -X dickobot/application/logging.GoVersion=$(go version | awk '{print $3}') \
        -X dickobot/application/logging.BuildRv=$(git describe --always --long)"
RUN chmod +x /app/dickobot

FROM busybox:1.36-musl

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=certs /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/dickobot /usr/local/bin/dickobot

RUN adduser -D -s /bin/sh dickobot
USER dickobot

CMD ["/usr/local/bin/dickobot"]