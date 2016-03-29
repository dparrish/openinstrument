package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument"
	"github.com/dparrish/openinstrument/client"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
)

var (
	connectAddress = flag.String("connect", "localhost:8021", "Connect directly to the specified datastore server.")
	readTimeout    = flag.Duration("read_timeout", 0, "Timeout for reading from stdin. Default is no timeout")
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	ctx := context.Background()

	conn, err := client.NewRpcClient(ctx, *connectAddress)
	if err != nil {
		log.Fatalf("Error connecting to %s: %s", *connectAddress, err)
	}
	defer conn.Close()

	re, err := regexp.Compile("^([^\\s]+) (.+)$")
	if err != nil {
		log.Fatalf("Error compiling regex: %s", err)
	}

	timestamp := openinstrument.NowMs()

	// Read from stdin
	c := make(chan *oproto.ValueStream)
	go func() {
		defer close(c)
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
				value.Value.(*oproto.Value_String_).String_ = matches[2]
			} else {
				value.Value.(*oproto.Value_Double).Double = f
			}
			c <- &oproto.ValueStream{
				Variable: v.AsProto(),
				Value:    []*oproto.Value{value},
			}
		}
	}()

	if *readTimeout != 0 {
		ctx, _ = context.WithTimeout(ctx, *readTimeout)
	}
	in, out, err := conn.Add(ctx)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		defer close(in)
		interval := time.Tick(1 * time.Second)
		request := &oproto.AddRequest{}
		for {
			select {
			case stream := <-c:
				if stream == nil {
					in <- request
					return
				}
				request.Stream = append(request.Stream, stream)
			case <-interval:
				in <- request
				request = &oproto.AddRequest{}
				timestamp = openinstrument.NowMs()
			case response := <-out:
				if response == nil {
					return
				}
				log.Printf("%s", response)
			}
		}
		in <- request
		log.Printf("done")
	}()

	for response := range out {
		log.Printf("%s", response)
	}
}
