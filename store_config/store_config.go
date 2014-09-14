package store_config

import (
  "code.google.com/p/goprotobuf/proto"
  "errors"
  "fmt"
  "log"
  "net"
  "os"
  "strings"
  oproto "github.com/dparrish/openinstrument/proto"
  "time"
)

type config struct {
  filename    string
  config_text string
  Config      *oproto.StoreConfig
  stat        os.FileInfo
}

var global_config *config

func NewConfig(filename string) (*config, error) {
  global_config = new(config)
  global_config.filename = filename
  if err := global_config.ReadConfig(); err != nil {
    return nil, err
  }
  go global_config.watchConfigFile()
  return global_config, nil
}

func Config() *config {
  if global_config == nil {
    log.Panic("global config not loaded")
  }
  return global_config
}

func (this *config) ReadConfig() error {
  file, err := os.Open(this.filename)
  if err != nil {
    return errors.New(fmt.Sprintf("Error opening config file %s: %s", this.filename, err))
  }
  defer file.Close()
  this.stat, _ = file.Stat()
  buf := make([]byte, this.stat.Size())
  n, err := file.Read(buf)
  if err != nil || int64(n) != this.stat.Size() {
    return errors.New(fmt.Sprintf("Error reading config file %s: %s", this.filename, err))
  }
  if err := this.handleNewConfig(string(buf)); err != nil {
    return err
  }
  return nil
}

func (this *config) handleNewConfig(text string) error {
  config := new(oproto.StoreConfig)
  err := proto.UnmarshalText(text, config)
  if err != nil {
    return errors.New(fmt.Sprintf("Error parsing config file %s: %s", this.filename, err))
  }
  this.Config = config
  return nil
}

func (this *config) watchConfigFile() {
  tick := time.Tick(10 * time.Second)
  for {
    select {
    case <-tick:
      stat, _ := os.Stat(this.filename)
      stat = stat
      this.ReadConfig()
    }
  }
}

func (this *config) ThisServer(port string) *oproto.StoreServer {
  addrs, _ := net.InterfaceAddrs()
  for _, server := range this.Config.Server {
    for _, addr := range addrs {
      parts := strings.SplitN(addr.String(), "/", 2)
      if server.GetAddress() == net.JoinHostPort(parts[0], port) {
        return server
      }
    }
  }
  return nil
}
