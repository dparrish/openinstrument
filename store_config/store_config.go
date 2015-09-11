package store_config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"code.google.com/p/goprotobuf/proto"
	oproto "github.com/dparrish/openinstrument/proto"
)

type config struct {
	filename   string
	configText string
	Config     *oproto.StoreConfig
	stat       os.FileInfo
}

var globalConfig *config

func NewConfig(filename string) (*config, error) {
	globalConfig = new(config)
	globalConfig.filename = filename
	if err := globalConfig.ReadConfig(); err != nil {
		return nil, err
	}
	go globalConfig.watchConfigFile()
	return globalConfig, nil
}

func Config() *config {
	if globalConfig == nil {
		log.Panic("global config not loaded")
	}
	return globalConfig
}

func (c *config) ReadConfig() error {
	file, err := os.Open(c.filename)
	if err != nil {
		return fmt.Errorf("Error opening config file %s: %s", c.filename, err)
	}
	defer file.Close()
	c.stat, _ = file.Stat()
	buf := make([]byte, c.stat.Size())
	n, err := file.Read(buf)
	if err != nil || int64(n) != c.stat.Size() {
		return fmt.Errorf("Error reading config file %s: %s", c.filename, err)
	}
	if err := c.handleNewConfig(string(buf)); err != nil {
		return err
	}
	return nil
}

func (c *config) handleNewConfig(text string) error {
	config := new(oproto.StoreConfig)
	err := proto.UnmarshalText(text, config)
	if err != nil {
		return fmt.Errorf("Error parsing config file %s: %s", c.filename, err)
	}
	c.Config = config
	return nil
}

func (c *config) watchConfigFile() {
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-tick:
			stat, _ := os.Stat(c.filename)
			stat = stat
			c.ReadConfig()
		}
	}
}

func (c *config) ThisServer(port string) *oproto.StoreServer {
	addrs, _ := net.InterfaceAddrs()
	for _, server := range c.Config.Server {
		for _, addr := range addrs {
			parts := strings.SplitN(addr.String(), "/", 2)
			if server.GetAddress() == net.JoinHostPort(parts[0], port) {
				return server
			}
		}
	}
	return nil
}
