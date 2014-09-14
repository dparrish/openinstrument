FROM golang:1.3
MAINTAINER David Parrish <david@dparrish.com>
ENV DEBIAN_FRONTEND noninteractive

# Download and build the store plus all dependencies
RUN apt-get -y update && apt-get -y install protobuf-compiler
RUN go get github.com/dparrish/openinstrument
#RUN go get github.com/coreos/go-etcd
#RUN go get github.com/joaojeronimo/go-crc16
#RUN go get github.com/nu7hatch/gouuid
RUN go get code.google.com/p/goprotobuf/proto
RUN go get code.google.com/p/goprotobuf/protoc-gen-go
RUN go install code.google.com/p/goprotobuf/protoc-gen-go
RUN mkdir /go/src/github.com/dparrish/openinstrument/proto
WORKDIR /go/src/github.com/dparrish/openinstrument
RUN protoc --go_out=proto --proto_path=. openinstrument.proto
WORKDIR /go
RUN go install github.com/dparrish/openinstrument/store

VOLUME /store
EXPOSE 8001
CMD ["-config", "/store/config.txt", "-datastore", "/store"]
ENTRYPOINT ["/go/bin/store"]

