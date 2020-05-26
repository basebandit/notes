package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/parish/notes/store/postgresql"

	"github.com/parish/notes/store/mysql"

	"github.com/parish/notes/config"
	"github.com/parish/notes/store"
)

const configFile = "notes.conf"

func initConfig() {
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		logger.Fatalf("configuration file already exists: %s", configFile)
	}

	logger.Printf("creating initial configuration: %s", configFile)

	cfg, err := config.Init()
	if err != nil {
		logger.Fatalf("failed to generate a configuration: %s", err)
	}

	err = ioutil.WriteFile(configFile, []byte(cfg), 0666)
	if err != nil {
		logger.Fatalf("failed to write configuration file: %s", err) 
	}
}

//loads configuration from environment variables or from file
func getConfig() (cfg *config.Config, err error) {
	if *useEnvConfig {
		cfg, err = config.ReadEnv()
	} else {
		cfg, err = config.ReadFile(configFile)
	}
	return cfg, err
}

//generates a random 32-byte hex-encoded Key.
func genKey() {
	logger.Printf("Key : %s", config.GenKey(32))
}

//returns a store handle (pointer) from the config
func getStore(cfg *config.Config) (store.Store, error) {
	switch cfg.Store.Type {
	case "mysql":
		return mysql.Connect(
			cfg.Store.MySQL.Address,
			cfg.Store.MySQL.Username,
			cfg.Store.MySQL.Password,
			cfg.Store.MySQL.Database,
		)
	case "postgresql":
		return postgresql.Connect(
			cfg.Store.PostrgeSQL.Address,
			cfg.Store.PostrgeSQL.Username,
			cfg.Store.PostrgeSQL.Password,
			cfg.Store.PostrgeSQL.Database,
			cfg.Store.PostrgeSQL.SSLMode,
			cfg.Store.PostrgeSQL.SSLRootCert,
		)
	}
	return nil, fmt.Errorf("unknown store type: %s", cfg.Store.Type)
}
