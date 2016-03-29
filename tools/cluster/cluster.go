package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dparrish/openinstrument"
	"github.com/dparrish/openinstrument/client"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/value"
	"github.com/dparrish/openinstrument/variable"
	"golang.org/x/net/context"
)

var (
	endKey         = flag.String("end_key", "", "End key")
	ID             = flag.String("id", "", "Block ID")
	connectAddress = flag.String("connect", "localhost:8021", "Connect directly to the specified datastore server.")
)

type empty struct{}
type semaphore chan empty

// acquire n resources
func (s semaphore) P(n int) {
	e := empty{}
	for i := 0; i < n; i++ {
		s <- e
	}
}

// release n resources
func (s semaphore) V(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}

func runSlowAddLoadtest(ctx context.Context) {
	log.Println("Running ADD load test")

	var sent, received int

	go func() {
		tick := time.Tick(1 * time.Second)
		for {
			select {
			case <-tick:
				log.Printf("Sent %d, received %d in the last second, latency %0.02fms", sent, received, 1.0/float64(received)*1000.0)
				received = 0
				sent = 0
			}
		}
	}()
	for {
		conn, err := client.NewRpcClient(ctx, *connectAddress)
		if err != nil {
			log.Fatalf("Error connecting to %s: %s", *connectAddress, err)
		}

		in, out, err := conn.Add(ctx)
		if err != nil {
			log.Printf("Error starting Add RPC: %s", err)
			return
		}

		v := variable.NewFromString("/test/var1{host=rage}").AsProto()
		request := &oproto.AddRequest{
			Stream: []*oproto.ValueStream{
				{
					Variable: v,
					Value: []*oproto.Value{
						value.NewDouble(openinstrument.NowMs(), 1.0),
					},
				},
			},
		}
		in <- request
		close(in)
		sent++
		<-out
		received++

		conn.Close()
	}
}

func runAddLoadtest(ctx context.Context, conn client.Client) {
	log.Println("Running ADD load test")
	in, out, err := conn.Add(ctx)
	if err != nil {
		log.Printf("Error starting Add RPC: %s", err)
		return
	}
	defer close(in)

	sem := make(semaphore, 1000)

	var sent, received int

	go func() {
		tick := time.Tick(1 * time.Second)
		for {
			select {
			case <-out:
				received++
				sem.V(1)
			case <-tick:
				log.Printf("Sent %d, received %d in the last second, latency %0.02fms", sent, received, 1.0/float64(received)*1000.0)
				received = 0
				sent = 0
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		sem.P(1)
		v := variable.NewFromString("/test/var1{host=rage}").AsProto()
		request := &oproto.AddRequest{
			Stream: []*oproto.ValueStream{
				{
					Variable: v,
					Value: []*oproto.Value{
						value.NewDouble(openinstrument.NowMs(), 1.0),
					},
				},
			},
		}
		in <- request
		sent++
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	ctx := context.Background()
	conn, err := client.NewRpcClient(ctx, *connectAddress)
	if err != nil {
		log.Fatalf("Error connecting to %s: %s", *connectAddress, err)
	}
	defer conn.Close()

	switch os.Args[1] {
	case "get":
		response, err := conn.GetConfig(ctx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(openinstrument.ProtoText(response))

	case "add_loadtest":
		runAddLoadtest(ctx, conn)
	case "slow_add_loadtest":
		runSlowAddLoadtest(ctx)
	}
}
