package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`

	Redis struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		DB   int    `yaml:"db"`
	} `yaml:"redis"`

	Kafka struct {
		KafkaURL string `yaml:"kafkaURL"`
		Topic    string `yaml:"topic"`
	} `yaml:"kafka"`

	Prometheus struct {
		Port int `yaml:"port"`
	} `yaml:"prometheus"`

	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
