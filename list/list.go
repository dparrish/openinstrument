package main

import (
  "code.google.com/p/goprotobuf/proto"
  "github.com/dparrish/openinstrument"
  openinstrument_proto "github.com/dparrish/openinstrument/proto"
  "github.com/dparrish/openinstrument/variable"
  "flag"
  "fmt"
  "log"
  "net/http"
  "os"
  "sort"
  "time"
)

var address = flag.String("address", "", "Address to listen on (blank for any)")
var port = flag.Int("port", 0, "Port to listen on")
var max_variables = flag.Int("max_variables", 0, "Maximum number of variables to return")
var max_age = flag.String("max_age", "", "Maximum history to search for old variables")
var config_file = flag.String("config", "/store/config.txt", "Path to the store configuration file")
var connect_address = flag.String("connect", "",
  "Connect directly to the specified datastore server. The Store config will be retrieved from this host and used.")

func List(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(w, "Hello")
}

// Argument server.
func Args(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintln(w, os.Args)
}

func main() {
  log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
  flag.Parse()
  if len(flag.Args()) < 1 {
    log.Fatal("Specify at least one variable to retrieve")
  }

  var client *openinstrument.StoreClient
  var err error
  if *connect_address != "" {
    client, err = openinstrument.NewAutoStoreClient(*connect_address)
    if err != nil {
      log.Fatal("Can't create StoreClient: %s", err)
    }
  } else if *config_file != "" {
    client, err = openinstrument.NewStoreClient(*config_file)
    if err != nil {
      log.Fatal("Can't create StoreClient: %s", err)
    }
  } else {
    log.Fatal("Specify either --connect or --config")
  }

  request := openinstrument_proto.ListRequest{
    Prefix: variable.NewFromString(flag.Arg(0)).AsProto(),
  }
  if *max_variables > 0 {
    request.MaxVariables = proto.Uint32(uint32(*max_variables))
  }
  if *max_age != "" {
    d, _ := time.ParseDuration(*max_age)
    request.MaxAge = proto.Uint64(uint64(d.Seconds() * 1000))
  }
  response, err := client.List(&request)
  if err != nil {
    log.Fatal(err)
  }
  vars := make([]string, 0)
  for _, listresponse := range response {
    for _, v := range listresponse.Variable {
      vars = append(vars, variable.NewFromProto(v).String())
    }
  }
  sort.Strings(vars)
  for _, variable := range vars {
    fmt.Println(variable)
  }
}