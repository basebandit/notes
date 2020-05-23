package config

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"

	"github.com/hashicorp/hcl"
)

//Config is our application's configuration struct(object).
type Config struct {
	Address string `hcl:"address" envconfig:"NOTE_APP_ADDRESS"`
	BaseURL string `hcl:"base_url" envconfig:"NOTE_APP_BASE_URL"`
	Title   string `hcl:"title" envconfig:"NOTE_APP_TITLE"`

	JWT struct {
		Secret string `hcl:"secret" envconfig:"NOTE_APP_JWT_SECRET"`
	} `hcl:"jwt"`

	Store struct {
		Type string `hcl:"type" envconfig:"NOTE_APP_STORE_TYPE"`

		PostrgeSQL struct {
			Address     string `hcl:"address" envconfig:"NOTE_APP_STORE_POSTGRESQL_ADDRESS"`
			Username    string `hcl:"username" envconfig:"NOTE_APP_STORE_POSTGRESQL_USERNAME"`
			Password    string `hcl:"password" envconfig:"NOTE_APP_STORE_POSTGRESQL_PASSWORD"`
			Database    string `hcl:"database" envconfig:"NOTE_APP_STORE_POSTGRESQL_DATABASE"`
			SSLMode     string `hcl:"sslmode" envconfig:"NOTE_APP_STORE_POSTGRESQL_SSLMODE"`
			SSLRootCert string `hcl:"sslrootcert" envconfig:"NOTE_APP_STORE_POSTGRESQL_SSLROOTCERT"`
		} `hcl:"postgresql"`

		MySQL struct {
			Address  string `hcl:"address" envconfig:"NOTE_APP_STORE_MYSQL_ADDRESS"`
			Username string `hcl:"username" envconfig:"NOTE_APP_STORE_MYSQL_USERNAME"`
			Password string `hcl:"password" envconfig:"NOTE_APP_STORE_MYSQL_PASSWORD"`
			Database string `hcl:"database" envconfig:"NOTE_APP_STORE_MYSQL_DATABASE"`
		} `hcl:"mysql"`
	} `hcl:"store"`
}

//ReadFile reads our application's config from file
func ReadFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	cfg := &Config{}
	err = hcl.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal hcl: %v", err)
	}

	prepare(cfg)
	return cfg, nil
}

//ReadEnv reads our application config from environment variables.
func ReadEnv() (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %v", err)
	}
	prepare(cfg)
	return cfg, nil
}

func prepare(cfg *Config) {
	cfg.BaseURL = strings.TrimSuffix(cfg.BaseURL, "/")
}

//Init generates an initial config string.
func Init() (string, error) {
	buf := new(bytes.Buffer)
	err := tpl.Execute(buf, map[string]interface{}{
		"jwt_secret": GenKey(32),
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

//GenKey generates a crypto-random key with byte length keyLen
//and hex-encodes it to a string.
func GenKey(keyLen int) string {
	bytes := make([]byte, keyLen)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(bytes)
}

var tpl = template.Must(template.New("initial-config").Parse(strings.TrimSpace(`
address  = "127.0.0.1:8080"
base_url = "https://notes.com/app"
title    = "notes"
jwt {
  secret = "{{.jwt_secret}}"
}
store {
  type = "postgresql"
  postgresql {
    address  = "127.0.0.1:5432"
    username = ""
    password = ""
    database = ""
    sslmode  = "disable"
  }
  mysql {
    address  = "127.0.0.1:3306"
    username = ""
    password = ""
    database = ""
  }
}
`)))
