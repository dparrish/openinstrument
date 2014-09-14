FROM golang:1.3
MAINTAINER David Parrish <david@dparrish.com>
ENV DEBIAN_FRONTEND noninteractive

# Download and build the store plus all dependencies
RUN go get github.com/dparrish/openinstrument/...
RUN go install github.com/dparrish/openinstrument/store

VOLUME /store
EXPOSE 8001
CMD ["-config", "/store/config.txt", "-datastore", "/store"]
ENTRYPOINT ["/go/bin/store"]

