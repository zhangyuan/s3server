FROM golang:1.19-alpine as builder

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build

FROM alpine:3.16.2

ARG USER=default
ENV HOME /home/$USER
RUN adduser -D $USER

COPY --from=builder /build/s3server /usr/local/bin/

USER $USER
WORKDIR $HOME

CMD [ "s3server" ]
