package config

import (
	"flag"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string `yaml:"env" env-default:"local"`
	GRPC GRPC   `yaml:"grpc"`
}

type GRPC struct {
	port    int           `yaml:"port"`
	timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchPathConfig()

	if path == "" {
		panic("config file path is empty")
	}
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("load config file failed: " + err.Error())
	}
	return &cfg
}

func fetchPathConfig() string {
	var path string

	flag.StringVar(&path, "config", "", "config file path")
	flag.Parse()

	if path == "" {
		panic("config file path is empty")
	}
	return path
}
