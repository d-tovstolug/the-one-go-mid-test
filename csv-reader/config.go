package main

import "github.com/caarlos0/env/v11"

type Config struct {
	Debug      bool   `env:"DEBUG" envDefault:"false"`
	WorkersNum int    `env:"WORKERS_NUM" envDefault:"10"`
	FileName   string `env:"FILENAME" envDefault:"./files/annual-enterprise-survey-2023-financial-year-provisional.csv"`
	SkipHeader bool   `env:"SKIP_HEADER" envDefault:"true"`
}

func ParseConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return &cfg, err
}
