FROM alpine:latest
MAINTAINER Patrick Debois <patrick.debois@jedi.be>

RUN apk -U add openssl
RUN apk --update add ca-certificates

COPY dist/ecs-watch_linux_amd64 /usr/local/bin/ecs-watch

ENTRYPOINT ["/usr/local/bin/ecs-watch"]
