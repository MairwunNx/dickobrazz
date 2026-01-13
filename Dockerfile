FROM alpine:3.23 AS certs
RUN apk --no-cache add ca-certificates tzdata

FROM golang:1.25.5-alpine3.21 AS builder

RUN apk --no-cache add git
COPY go.mod go.sum ./
RUN go mod download

COPY program.go ./
COPY .version ./
COPY .git ./
COPY application/ ./application/

RUN mkdir -p /app
RUN go build -o /app/dickobrazz \
    -a -installsuffix cgo \
    -gcflags "all=-N -l" \
    -tags timetzdata \
    -ldflags="-s -w \
        -X dickobrazz/application/logging.Version=$(cat .version) \
        -X dickobrazz/application/logging.BuildAt=$(date +%Y-%m-%d_%H:%M:%S) \
        -X dickobrazz/application/logging.GoVersion=$(go version | awk '{print $3}') \
        -X dickobrazz/application/logging.BuildRv=$(git describe --always --long)"
RUN chmod +x /app/dickobrazz

FROM busybox:1.37-musl

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=certs /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/dickobrazz /usr/local/bin/dickobrazz

RUN adduser -D -s /bin/sh dickobrazz
USER dickobrazz

EXPOSE 80

CMD ["/usr/local/bin/dickobrazz"]