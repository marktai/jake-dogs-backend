package main

import (
	"corgi"
	"flag"
	"server"
	"time"
)

func main() {

	var minutes int
	var port int

	var max_images int
	var max_gifs int

	flag.IntVar(&minutes, "minutes", 60, "Period of checking")
	flag.IntVar(&port, "port", 8014, "Port that the server listens on")

	flag.IntVar(&max_images, "max-images", 500, "Number of images that are stored at max")
	flag.IntVar(&max_gifs, "max-gifs", 500, "Number of gifs that are stored at max")

	flag.Parse()

	timeout := time.Duration(minutes) * time.Minute

	corgi.Init(max_images, max_gifs, &timeout)
	server.Run(port)
}
