FROM alpine:edge

ENV GOPATH /go
ENV GOREPO github.com/allen13/influxdb-write-simulator
RUN mkdir -p $GOPATH/src/$GOREPO
COPY . $GOPATH/src/$GOREPO
WORKDIR $GOPATH/src/$GOREPO

RUN set -ex \
	&& apk add --no-cache --virtual .build-deps \
		git \
		go \
		build-base \
	&& go build influxdb-write-simulator.go \
	&& apk del .build-deps \
	&& rm -rf $GOPATH/pkg

ENV WRITE_INTERVAL 10s
ENV INFLUXDB_HOST influxdb
ENV INFLUXDB_PORT 8086
ENV INFLUXDB_USER user
ENV INFLUXDB_PASSWD password


CMD ./influxdb-write-simulator
