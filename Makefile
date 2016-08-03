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
	go get github.com/olekukonko/tablewriter

BUILD_VERSION ?= 1
BUILD_DATE=now

dist:
	gox -ldflags "-X main.Version $(BUILD_VERSION) -X main.BuildDate $(BUILD_DATE)" -output "dist/ecs-watch_{{.OS}}_{{.Arch}}"

release:
	ghr -t $(GITHUB_TOKEN) -u jedi4ever -r ecs-watch --replace `git describe --tags` dist/
