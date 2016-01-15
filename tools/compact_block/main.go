package main

import (
	"flag"
	"log"
	"os"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	endKey         = flag.String("end_key", "", "End key")
	ID             = flag.String("id", "", "Block ID")
	connectAddress = flag.String("connect", "localhost:8021", "Connect directly to the specified datastore server.")
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	conn, err := grpc.Dial(*connectAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting to %s: %s", *connectAddress, err)
	}
	defer conn.Close()

	request := &oproto.CompactBlockRequest{
		Block: &oproto.Block{},
	}
	if *endKey != "" {
		request.Block.EndKey = *endKey
	} else if *ID != "" {
		request.Block.Id = *ID
	} else if len(os.Args) > 1 {
		request.Block.Id = os.Args[1]
	} else {
		log.Fatal("Specify either --end_key or --id")
	}

	stub := oproto.NewStoreClient(conn)
	response, err := stub.CompactBlock(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(openinstrument.ProtoText(response))
}
