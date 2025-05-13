FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /feedjson ./cmd/feedjson

FROM alpine:3.18


RUN apk add --no-cache ca-certificates

COPY --from=builder /feedjson /feedjson

EXPOSE 8080

ENTRYPOINT ["/feedjson"]
