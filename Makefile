.PHONY: build report get

STS_EXEC ?= 

ecs-watch: *.go
	go build

build: ecs-watch
report: build
	$(STS_EXEC) ./ecs-watch report

get:
	go get github.com/tj/go-debug
	go get github.com/urfave/cli
