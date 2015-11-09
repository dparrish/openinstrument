all: proto/openinstrument.pb.go query/lexer/lexer.go

proto/openinstrument.pb.go: openinstrument.proto
	PATH=${PATH}:${GOPATH}/bin protoc --go_out=plugins=grpc:proto --proto_path=. $<

test:
	go test github.com/dparrish/openinstrument/...

query/lexer/lexer.go: query.bnf
	${GOPATH}/bin/gocc -o query -p query $<
	find query -name "*.go" -exec chmod -x {} \;
