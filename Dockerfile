FROM golang:alpine as build

RUN apk update && apk add curl git

RUN go get golang.org/x/sys/unix \
    && go get golang.org/x/crypto/ssh \
    && go get gopkg.in/ini.v1 \
    && go get github.com/julienschmidt/httprouter

WORKDIR /go/src/github.com

COPY ./gw_stats_iox ./gw_stats_iox

WORKDIR /go/src/github.com/gw_stats_iox

RUN go build -o gw_server gw_server.go

FROM alpine

COPY --from=build /go/src/github.com/gw_stats_iox/gw_server .
COPY ./gw_stats_iox/package_config.ini .

EXPOSE 8080

CMD ["./gw_server"]



