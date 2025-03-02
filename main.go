package main

import (
	"flag"

	"github.com/DisguisedTrolley/chirpy/api"
	"github.com/charmbracelet/log"
	_ "github.com/lib/pq"
)

const (
	PORT = "8080"
)

func main() {
	var port string
	flag.StringVar(&port, "p", PORT, "Port to run the server on")

	server := api.NewServer(port)
	config := api.NewApiConfig()

	err := server.Start(config)
	if err != nil {
		log.Errorf("Error: Unable to start server, %s\n", err.Error())
	}
}
