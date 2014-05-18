package sup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// All configuration for a Scramble server+notary.
// The config object is read from ~/.sup/config.json
type Config struct {
	DbServer   string
	DbUser     string
	DbPassword string
	DbCatalog  string

	HTTPPort int // internal, nginx handles SSL and forwards
}

// Gets the cotents of the Scramble config file, ~/.sup/config.json
// The file is read only once at startup.
func GetConfig() *Config {
	return &config
}

func validateConfig(cfg *Config) error {
	if cfg.DbServer == "" {
		return errors.New("DbServer must be set")
	}
	if cfg.DbUser == "" {
		return errors.New("DbUser must be set")
	}
	if cfg.DbCatalog == "" {
		return errors.New("DbCatalog must be set")
	}
	if cfg.HTTPPort == 0 {
		return errors.New("HTTPPort must be set")
	}
	return nil
}

var defaultConfig = Config{
	"127.0.0.1",
	"sup",
	"sup",
	"sup",

	8888,
}

var config Config

func init() {
	configFile := os.Getenv("HOME") + "/.sup/config.json"
	log.Printf("Reading " + configFile)

	// try to read configuration. if missing, write default
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		writeDefaultConfig(configFile)
		fmt.Println("Config file written to ~/.sup/config.json. Please edit & run again")
		os.Exit(1)
		return
	}

	// try to parse configuration. on error, die
	config = Config{}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Panicf("Invalid configuration file %s: %v", configFile, err)
	}
	err = validateConfig(&config)
	if err != nil {
		log.Panicf("Invalid configuration file %s: %v", configFile, err)
	}
}

func writeDefaultConfig(configFile string) {
	log.Printf("Creating default configration file %s", configFile)
	configBytes, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(filepath.Dir(configFile), 0700)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(configFile, configBytes, 0600)
	if err != nil {
		panic(err)
	}
}

