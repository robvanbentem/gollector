package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type config struct {
	MQTT configMQTT

	LogFile string

	Database string
	DBUser   string
	DBPass   string
	DBHost   string
	DBPort   string
	DBScheme string

	Source string
}

type configMQTT struct {
	Broker string
	Name   string
	QOS    byte
	Path   string
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

var ConfigRoot *config

func LoadConfig() {
	buf, err := ioutil.ReadFile("Config.toml")
	if err != nil {
		panic("Could not read Config.toml")
	}

	var conf config
	if _, err := toml.Decode(string(buf), &conf); err != nil {
		panic("Could not parse Config.toml: " + err.Error())
	}

	var setting string
	var found bool

	if setting, found = os.LookupEnv("DB_STR"); found == true {
		conf.Database = setting
	}

	if setting, found = os.LookupEnv("DB_USER"); found == true {
		conf.DBUser = setting
	}

	if setting, found = os.LookupEnv("DB_PASS"); found == true {
		conf.DBPass = setting
	}

	if setting, found = os.LookupEnv("DB_HOST"); found == true {
		conf.DBHost = setting
	}

	if setting, found = os.LookupEnv("DB_PORT"); found == true {
		conf.DBPort = setting
	}

	if setting, found = os.LookupEnv("DB_SCHEME"); found == true {
		conf.DBScheme = setting
	}

	if setting, found = os.LookupEnv("APP_SOURCE"); found == true {
		conf.Source = setting
	}

	if conf.DBScheme != "" {
		conf.Database = fmt.Sprintf(conf.Database, conf.DBUser, conf.DBPass, conf.DBHost, conf.DBPort, conf.DBScheme)
	}

	ConfigRoot = &conf
}
