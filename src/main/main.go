package main

import (
	"flag"
	"subredditCrawler"
	"time"
)

func main() {

	var port int

	flag.IntVar(&port, "Port", 8080, "Port the server listens to")

	flag.Parse()

	killChan := make(chan bool)

	subredditCrawler.Run("buildapcsales", "G502", 30*time.Minute, killChan)
}
