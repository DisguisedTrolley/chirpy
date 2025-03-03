package main

import (
	"flag"

	"github.com/DisguisedTrolley/chirpy/app/api"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	PORT = "8080"
)

func main() {
	godotenv.Load("../.env")

	var port string
	flag.StringVar(&port, "p", PORT, "Port to run the server on")

	server := api.NewServer(port)
	config := api.NewApiConfig()

	err := server.Start(config)
	if err != nil {
		log.Errorf("Error: Unable to start server, %s\n", err.Error())
	}
}
