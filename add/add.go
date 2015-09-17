package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strconv"

	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"

	"google.golang.org/grpc"
)

var (
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

	re, err := regexp.Compile("^([^\\s]+) (.+)$")
	if err != nil {
		log.Fatalf("Error compiling regex: %s", err)
	}

	request := &oproto.AddRequest{}
	timestamp := openinstrument.NowMs()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		matches := re.FindStringSubmatch(scanner.Text())
		if len(matches) != 3 {
			log.Printf("Invalid input line: %s", scanner.Text())
			continue
		}
		v := variable.NewFromString(matches[1])
		value := &oproto.Value{Timestamp: timestamp}
		f, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			value.StringValue = matches[2]
		} else {
			value.DoubleValue = f
		}
		stream := &oproto.ValueStream{
			Variable: v.AsProto(),
			Value:    []*oproto.Value{value},
		}
		request.Stream = append(request.Stream, stream)
	}
	log.Printf("%s", request)

	if len(request.Stream) == 0 {
		log.Fatalf("No valid values to add")
	}

	stub := oproto.NewStoreClient(conn)
	response, err := stub.Add(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", response)
}
