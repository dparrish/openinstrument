package main

import (
	"flag"
	"fmt"
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

	switch os.Args[1] {
	case "get":
		request := &oproto.GetClusterRequest{}
		stub := oproto.NewStoreClient(conn)
		response, err := stub.GetCluster(context.Background(), request)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(openinstrument.ProtoText(response))
	}
}
