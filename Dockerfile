FROM golang:1.19-alpine as builder

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build

FROM alpine:3.16.2

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN apk add --no-cache tini

ARG USER=default
ENV HOME /home/$USER
RUN adduser -D $USER

COPY --from=builder /build/s3server /usr/local/bin/

USER $USER
WORKDIR $HOME

ENTRYPOINT ["/sbin/tini", "--"]
CMD [ "s3server" ]
