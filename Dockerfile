FROM golang:1.13.6-alpine3.11  as builder
RUN apk add --update git alpine-sdk

ENV GOPATH=/go

COPY . $GOPATH/src/auth/
WORKDIR $GOPATH/src/auth/
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
RUN go build -o auth .

FROM alpine

COPY --from=builder /go/src/auth/auth /auth

EXPOSE  8080

ENTRYPOINT ["/auth","migrate","up"]