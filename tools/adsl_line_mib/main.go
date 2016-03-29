package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument"
	"github.com/dparrish/openinstrument/client"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/value"
	"github.com/dparrish/openinstrument/variable"
	"github.com/soniah/gosnmp"
)

var (
	connectAddress = flag.String("connect", "localhost:8021", "Connect directly to the specified datastore server.")
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
	in, out, err := conn.Add(ctx)
	if err != nil {
		log.Fatal(err)
	}

	snmp := gosnmp.Default
	for _, host := range []string{"192.168.1.1"} {
		snmp.Target = host
		snmp.Community = "public"
		if err := snmp.Connect(); err != nil {
			log.Fatal(err)
		}
		defer snmp.Conn.Close()

		// VDSL-LINE-MIB::vdslPhysCurrLineRate.1 = Gauge32: 20265 kbps
		// VDSL-LINE-MIB::vdslPhysCurrLineRate.2 = Gauge32: 1315 kbps
		oids := []string{"1.3.6.1.2.1.10.97.1.1.2.1.10.1", "1.3.6.1.2.1.10.97.1.1.2.1.10.2"}
		result, err := snmp.Get(oids)
		if err != nil {
			log.Fatal(err)
		}

		request := &oproto.AddRequest{
			Stream: []*oproto.ValueStream{},
		}

		for _, snmpVar := range result.Variables {
			switch snmpVar.Name {
			case ".1.3.6.1.2.1.10.97.1.1.2.1.10.1":
				request.Stream = append(request.Stream, &oproto.ValueStream{
					Variable: variable.NewFromString(fmt.Sprintf("/network/adsl/downstream_rate{host=%s,unit=bps}", host)).AsProto(),
					Value: []*oproto.Value{
						value.NewDouble(openinstrument.NowMs(), float64(snmpVar.Value.(uint)*1000)),
					},
				})
			case ".1.3.6.1.2.1.10.97.1.1.2.1.10.2":
				request.Stream = append(request.Stream, &oproto.ValueStream{
					Variable: variable.NewFromString(fmt.Sprintf("/network/adsl/upstream_rate{host=%s,unit=bps}", host)).AsProto(),
					Value: []*oproto.Value{
						value.NewDouble(openinstrument.NowMs(), float64(snmpVar.Value.(uint)*1000)),
					},
				})
			}
		}

		in <- request
	}

	close(in)
	for response := range out {
		if !response.Success {
			fmt.Println(openinstrument.ProtoText(response))
		}
	}
}
