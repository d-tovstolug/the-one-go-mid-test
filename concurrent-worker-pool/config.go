package main

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Debug                 bool `env:"DEBUG" envDefault:"true"`
	WorkersNum            int  `env:"WORKERS_NUM" envDefault:"3"`
	TasksNum              int  `env:"TASKS_NUM" envDefault:"50"`
	MinWaitTimeSeconds    int  `env:"MIN_WAIT_TIME_SECONDS" envDefault:"1"`
	MaxWaitTimeSeconds    int  `env:"MAX_WAIT_TIME_SECONDS" envDefault:"3"`
	ErrProbabilityPercent int  `env:"ERROR_PROBABILITY_PERCENT" envDefault:"10"`
}

func ParseConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return &cfg, err
}
