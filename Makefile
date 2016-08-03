.PHONY: build report get generate dist dockerdist dockerbuild release

STS_EXEC ?= 

ecs-watch: *.go
	go build

build: ecs-watch

report: build
	$(STS_EXEC) ./ecs-watch report

generate: build
	$(STS_EXEC) ./ecs-watch generate --template-file nginx.tmpl

get:
	go get github.com/tj/go-debug
	go get github.com/urfave/cli
	go get github.com/olekukonko/tablewriter

BUILD_VERSION ?= 1
BUILD_DATE=now

dist:
	gox -ldflags "-X main.Version $(BUILD_VERSION) -X main.BuildDate $(BUILD_DATE)" -output "dist/ecs-watch_{{.OS}}_{{.Arch}}"

dockerdist: ecs-watch
	gox -osarch="linux/amd64" -ldflags "-X main.Version $(BUILD_VERSION) -X main.BuildDate $(BUILD_DATE)" -output "dist/ecs-watch_{{.OS}}_{{.Arch}}"

dockerbuild: dockerdist
	docker build  . -t ecs-watch:develop

release:
	ghr -t $(GITHUB_TOKEN) -u jedi4ever -r ecs-watch --replace `git describe --tags` dist/
