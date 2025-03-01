package main

import (
	"fmt"
	"net/http"
)

const (
	PORT = "8080"
	URL  = "http://localhost"
)

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    URL + ":" + PORT,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error: Unable to start server, %s\n", err.Error())
	}
}
