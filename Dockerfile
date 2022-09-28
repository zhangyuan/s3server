FROM golang:1.19-alpine as builder

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build

FROM alpine:3.16.2

WORKDIR /

COPY --from=builder /build/s3server .

CMD [ "s3server" ]
