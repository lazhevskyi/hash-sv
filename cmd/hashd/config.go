package main

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.Bool("debug", true, "")
	pflag.Duration("hash_ttl", 5*time.Second, "time to regenerate hash")
	pflag.Int("http_port", 8080, "http server port")
	pflag.Int("grpc_port", 8081, "grpc server port")
}

type config struct {
	Debug    bool          `mapstructure:"debug"`
	HashTTL  time.Duration `mapstructure:"hash_ttl"`
	HttpPort int           `mapstructure:"http_port"`
	GrpcPort int           `mapstructure:"grpc_port"`
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
