package main

import (
	"flag"
	"log"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	varName        = flag.String("variable", "", "Variable")
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

	request := &oproto.LookupBlockRequest{}
	if *varName != "" {
		request.Variable = variable.NewFromString(*varName).AsProto()
	} else if *ID != "" {
		request.BlockId = *ID
	} else {
		log.Fatal("Specify either --variable or --id")
	}

	stub := oproto.NewStoreClient(conn)
	response, err := stub.LookupBlock(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(openinstrument.ProtoText(response))
}
