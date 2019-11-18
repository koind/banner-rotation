package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"time"
)

// The path to the configuration
var Path string

// Settings microservice
type Options struct {
	Postgres   Postgres
	GRPCServer GRPCServer
	HTTPServer HTTPServer
}

// Initializes microservice configurations
func Init(configPath string) Options {
	opt := Options{}

	if _, err := toml.DecodeFile(configPath, &opt); err != nil {
		log.Fatal("Не удалось загрузить конфиги микросервиса ", err)
	}

	return opt
}

// Settings postgres
type Postgres struct {
	DSN             string
	PingTimeout     int
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// Settings gRPC server
type GRPCServer struct {
	Host string
	Port int
}

// Returns the domain
func (s GRPCServer) GetDomain() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// Settings HTTP server
type HTTPServer struct {
	Host string
	Port int
}

// Returns the domain
func (s HTTPServer) GetDomain() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
