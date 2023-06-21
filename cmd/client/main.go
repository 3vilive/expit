package main

import (
	"flag"
	"log"
	"net"

	"github.com/3vilive/expit/cmd/client/internal/pkg/config"
	"github.com/3vilive/expit/cmd/client/internal/pkg/tunnel"
)

func init() {
	flag.StringVar(&config.ServerAddr, "server", "127.0.0.1:20000", "")
	flag.StringVar(&config.LocalAppAddr, "local", "127.0.0.1:8000", "")
}

func main() {
	flag.Parse()

	log.Printf("server: %s, local: %s\n", config.ServerAddr, config.LocalAppAddr)

	// check arguments
	if _, err := net.ResolveTCPAddr("tcp", config.ServerAddr); err != nil {
		log.Fatalf("resolve server address error: %s\n", err)
	}
	if _, err := net.ResolveTCPAddr("tcp", config.LocalAppAddr); err != nil {
		log.Fatalf("resolve local address error: %s\n", err)
	}

	// start tunnel
	controlTunnel, err := tunnel.NewControlTunnel(config.ServerAddr)
	if err != nil {
		log.Fatalf("failed to create control tunnel: %s\n", err)
	}
	controlTunnel.Start()
}
