package main

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.Bool("debug", true, "")
	pflag.Duration("hash_ttl", 5*time.Second, "time to regenerate hash")
	pflag.String("port", "8080", "http server port")
	pflag.String("addr", "", "http server addr")
}

type config struct {
	Debug   bool          `mapstructure:"debug"`
	HashTTL time.Duration `mapstructure:"hash_ttl"`
	Port    string        `mapstructure:"port"`
	Addr    string        `mapstructure:"addr"`
}

func mustParseConfig() config {
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(err)
	}

	var cfg config
	err = viper.UnmarshalExact(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
