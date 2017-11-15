export GOPATH := $(shell pwd)
default: build

build: 
	docker build -t jake-dogs . 

run: build
	-docker kill $(docker ps -q --filter ancestor="jake-dogs:latest")
	docker run -p 8013:8013 -t jake-dogs:latest > log.txt 2>&1 &

log: run
	tail -f -n2 log.txt
