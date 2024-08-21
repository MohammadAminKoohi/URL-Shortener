FROM golang:1.22.5-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o urlshortener ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/urlshortener .

EXPOSE 1323

CMD ["./urlshortener"]