package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/golang/protobuf/proto"
)

var (
	connectAddress = flag.String("connect", "localhost:8021", "Connect directly to the specified datastore server.")
	readTimeout    = flag.Duration("read_timeout", 0, "Timeout for reading from stdin. Default is no timeout")
)

func readAndParsePolicy() *oproto.RetentionPolicyItem {
	item := &oproto.RetentionPolicyItem{}
	policyTxt, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading stdin: %s", err)
	}

	if err := proto.UnmarshalText(string(policyTxt), item); err != nil {
		log.Fatalf("Invalid policy text: %s", err)
	}
	return item
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	ctx := context.Background()

	conn, err := grpc.Dial(*connectAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting to %s: %s", *connectAddress, err)
	}
	defer conn.Close()
	stub := oproto.NewStoreClient(conn)

	request := &oproto.UpdateRetentionPolicyRequest{
		Item: &oproto.RetentionPolicyItem{},
	}

	if len(os.Args) < 2 {
		log.Fatalf("Invalid command (append/insert/delete)")
	}

	switch os.Args[1] {
	case "append":
		request.Op = oproto.UpdateRetentionPolicyRequest_APPEND
		request.Item = readAndParsePolicy()
	case "insert":
		if len(os.Args) < 3 {
			log.Fatalf("Specify position (0-indexed)")
		}
		position, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid position: %s", err)
		}
		request.Op = oproto.UpdateRetentionPolicyRequest_INSERT
		request.Position = uint32(position)
		request.Item = readAndParsePolicy()
	case "delete":
		if len(os.Args) < 3 {
			log.Fatalf("Specify position (0-indexed)")
		}

		position, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid position: %s", err)
		}
		request.Op = oproto.UpdateRetentionPolicyRequest_REMOVE
		request.Position = uint32(position)
	default:
		log.Fatalf("Invalid command (append/insert/delete)")
	}

	response, err := stub.UpdateRetentionPolicy(ctx, request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(openinstrument.ProtoText(response))
}
