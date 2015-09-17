all: proto/openinstrument.pb.go

proto/openinstrument.pb.go: openinstrument.proto
	PATH=${PATH}:${GOPATH}/bin protoc --go_out=plugins=grpc:proto --proto_path=. openinstrument.proto

test:
	go test github.com/dparrish/openinstrument/...
