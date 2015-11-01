package config

import (
	"fmt"
	toml "github.com/BurntSushi/toml"
	"io/ioutil"
	// "os"
	"database/sql"
	"strconv"
	// _ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/godrv"
)

var Config Configuration

// the types

type Configuration struct {
	Debug       int `toml:"debug"`
	Fail        bool
	Server      serverConfig   `toml:"server"`
	PDF         pdfConfig      `toml:"pdf"`
	Database    databaseConfig `toml:"db"`
	OldDatabase databaseConfig `toml:"old_db"`
	Github      githubConfig   `toml:"github"`
	Partial		bool           `toml:"partial"`
}

type serverConfig struct {
	Hostname string `toml:"hostname"`
	Port     int    `toml:"port"`
}

type pdfConfig struct {
	Path string `toml:"path"`
}

type databaseConfig struct {
	Hostname string `toml:"host"`
	Database string `toml:"db"`
	Username string `toml:"user"`
	Password string `toml:"password"`
}

type githubConfig struct {
	AccessToken string `toml:"access_token"`
}

// loading the config

func init() {
	LoadConfig(true)
}

func LoadConfig(initial bool) {
	config := Configuration{
		Debug: 0,
		Fail:  false,
		Server: serverConfig{
			Hostname: "",
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

	config.Fail = false

	// if that worked, swap the config for the new one
	Config = config
	if Config.Debug > 0 {
		DebugConfig()
	}
}

func DebugConfig() {
	fmt.Printf("Config: %#v\n", Config)
}

func (server *serverConfig) Host() string {
	return server.Hostname + ":" + strconv.Itoa(server.Port)
}

func (db *databaseConfig) Open() (*sql.DB, error) {
	conn := db.Database + "/" + db.Username + "/" + db.Password
	if db.Hostname != "localhost" && db.Hostname != "" {
		conn = "tcp:" + db.Hostname + "*" + conn
	}
	if Config.Debug > 0 {
		fmt.Println("Connecting to", conn)
	}
	sqldb, err := sql.Open("mymysql", conn)
	return sqldb, err
}
