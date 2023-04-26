package config

import (
	"github.com/caarlos0/env/v8"
)

type Config struct {
	UIURL      string `env:"UI_URL"`
	S3Endpoint string `env:"S3_ENDPOINT"`
	S3Bucket   string `env:"S3_BUCKET"`
	S3KeyId    string `env:"S3_ID"`
	S3Secret   string `env:"S3_SECRET"`
}

func FromEnv() Config {
	configOpts := env.Options{
		Prefix:          "BUDDAFEST_REGISTRATION_",
		RequiredIfNoDef: true,
	}

	var cfg Config

	if err := env.ParseWithOptions(&cfg, configOpts); err != nil {
		panic(err)
	}

	return cfg
}
