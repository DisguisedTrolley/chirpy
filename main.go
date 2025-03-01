package main

import (
	"flag"

	"github.com/DisguisedTrolley/chirpy/api"
	"github.com/charmbracelet/log"
)

const (
	PORT = "8080"
)

func main() {
	var port string
	flag.StringVar(&port, "p", PORT, "Port to run the server on")

	server := api.NewServer(port)

	err := server.Start()
	if err != nil {
		log.Errorf("Error: Unable to start server, %s\n", err.Error())
	}
}
