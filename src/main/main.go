package main

import (
	"jake_dogs"
	"flag"
	"server"
)

func main() {
	var port int

	flag.IntVar(&port, "port", 8013, "Port that the server listens on")

	flag.Parse()

	jake_dogs.Init()
	server.Run(port)
}
