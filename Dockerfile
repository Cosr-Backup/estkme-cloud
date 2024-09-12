# Build: estkme-cloud
FROM golang:1.23-alpine AS builder

WORKDIR /app

ARG VERSION

COPY . .

RUN apk add --no-cache gcc musl-dev

RUN set -eux \
    && CGO_ENABLED=1 go build -trimpath -ldflags="-w -s -X main.Version=${VERSION}" -o estkme-cloud main.go

# Production
FROM alpine:3.20 AS production

WORKDIR /app

COPY --from=builder /app/estkme-cloud /app/estkme-cloud

RUN set -eux \
    && apk add --no-cache libcurl

EXPOSE 1888

CMD ["/app/estkme-cloud"]
