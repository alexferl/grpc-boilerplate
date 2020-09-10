package grpcd

import (
	"net"

	"github.com/spf13/pflag"

	"grpc-boilerplate/internal/pkg/config"
)

// Config holds all configuration for our program
type Config struct {
	config.Config
	BindAddress       net.IP
	BindPort          uint
	GracefulTimeout   uint
	HealthBindAddress net.IP
	HealthBindPort    uint
}

// NewConfig creates a Config instance
func NewConfig() Config {
	cnf := Config{
		Config:            config.NewConfig(),
		BindAddress:       net.ParseIP("127.0.0.1"),
		BindPort:          50051,
		GracefulTimeout:   30,
		HealthBindAddress: net.ParseIP("127.0.0.1"),
		HealthBindPort:    50052,
	}
	return cnf
}

// addFlags adds all the flags from the command line
func (cnf *Config) addFlags(fs *pflag.FlagSet) {
	fs.IPVar(&cnf.BindAddress, "bind-address", cnf.BindAddress, "The IP address to listen at.")
	fs.UintVar(&cnf.BindPort, "bind-port", cnf.BindPort, "The port to listen at.")
	fs.UintVar(&cnf.GracefulTimeout, "graceful-timeout", cnf.GracefulTimeout,
		"Timeout for graceful shutdown.")
	fs.IPVar(&cnf.HealthBindAddress, "health-bind-address", cnf.HealthBindAddress, "The IP address to listen at.")
	fs.UintVar(&cnf.HealthBindPort, "health-bind-port", cnf.HealthBindPort, "The port to listen at.")
}

// BindFlags normalizes and parses the command line flags
func (cnf *Config) BindFlags() {
	cnf.addFlags(pflag.CommandLine)
	cnf.Config.BindFlags()
}
