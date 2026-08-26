FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o auth .

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

COPY --from=builder /build/auth /auth

EXPOSE 8080

ENTRYPOINT ["/auth"]
