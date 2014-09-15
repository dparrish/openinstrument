FROM golang:1.3
MAINTAINER David Parrish <david@dparrish.com>

# Download and build the store plus all dependencies
RUN go get github.com/dparrish/openinstrument/...
RUN go install github.com/dparrish/openinstrument/store

VOLUME /store
EXPOSE 8020
CMD ["-config", "/store/config.txt"]
ENTRYPOINT ["/go/bin/store", "-templates", "src/github.com/dparrish/openinstrument/html"]

