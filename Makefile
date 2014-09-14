all: proto/openinstrument.pb.go

proto/openinstrument.pb.go: openinstrument.proto
	protoc --go_out=proto --proto_path=. openinstrument.proto
