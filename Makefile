export GOPATH := $(shell pwd)
default: build

init:
	rm -f bin/main bin/subredCrawler 
	@cd src/main && go get

build: init
	go build -o bin/subredCrawler src/main/main.go 

run: build
	@-pkill subredCrawler
	bin/subredCrawler >>log.txt 2>&1 &
