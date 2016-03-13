package main

import (
	"flag"
	"subredditCrawler"
	"time"
)

func main() {

	var minutes int
	var subreddit string
	var exp string
	var email string

	flag.IntVar(&minutes, "minutes", 30, "Period of checking")
	flag.StringVar(&subreddit, "subreddit", "buildapcsales", "Subreddit to check")
	flag.StringVar(&exp, "expression", "G502", "Non-case-sensitive expression to check for")
	flag.StringVar(&email, "email", "taifighterm@gmail.com", "Email to send the notification to")

	flag.Parse()

	killChan := make(chan bool)

	subredditCrawler.Run(subreddit, exp, email, time.Duration(minutes)*time.Minute, killChan)
}
