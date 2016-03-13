export GOPATH := $(shell pwd)
default: build

init:
	rm -f bin/main bin/subredditCrawler 
	@cd src/main && go get

build: init
	go build -o bin/subredditCrawler src/main/main.go 

run: build
	@-pkill subredditCrawler
	bin/subredditCrawler >>log.txt 2>&1 &
