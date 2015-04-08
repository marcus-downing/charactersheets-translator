package config

import (
	"fmt"
	toml "github.com/BurntSushi/toml"
	"io/ioutil"
	// "os"
	"strconv"
)

var Config Configuration

// the types

type Configuration struct {
	Debug    int `toml:"debug"`
	Fail     bool
	Server   ServerConfig
	PDF      PDFConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Hostname string `toml:"hostname"`
	Port     int    `toml:"port"`
}

type PDFConfig struct {
	Path string `toml:"path"`
}

type DatabaseConfig struct {
	Hostname string `toml:"hostname"`
}

// loading the config

func init() {
	LoadConfig(true)
}

func LoadConfig(initial bool) {
	config := Configuration{
		Debug: 1,
		Fail:  false,
		Server: ServerConfig{
			Hostname: "localhost",
			Port:     9091,
		},
	}
	if initial {
		Config = config
	}

	configData, err := ioutil.ReadFile("config.toml")
	if err != nil {
		fmt.Println("Error opening config.toml:", err)
		Config.Fail = true
		return
	}
	if _, err := toml.Decode(string(configData), &config); err != nil {
		// handle error
		fmt.Println("Error decoding config.toml:", err)
		Config.Fail = true
		return
	}

	if Config.Debug > 0 {
		DebugConfig()
	}
	// if that worked, swap the config for the new one
	Config = config
}

func DebugConfig() {
	fmt.Printf("Config: %#v\n", Config)
}

func (server ServerConfig) Host() string {
	return server.Hostname + ":" + strconv.Itoa(server.Port)
}
