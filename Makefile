export GOPATH := $(shell pwd)
default: build

init:
	rm -f bin/main bin/corgi-server 
	@cd src/main && go get

build: init
	go build -o bin/corgi-server src/main/main.go 

run: build
	@-pkill corgi-server
	bin/corgi-server >>log.txt 2>&1 &

log: run
	tail -f -n2 log.txt
