FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .env

ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o gowebly_fiber

RUN apk --no-cache add ca-certificates

FROM scratch

COPY --from=builder /build/gowebly_fiber /
COPY --from=builder /build/static /static
COPY --from=builder /build/.env .env

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/gowebly_fiber"]
