package main

import (
	//"errors"
	//"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/urfave/cli"
	"os"
)

import godebug "github.com/tj/go-debug"

var debug = godebug.Debug("ecs-watch")

var clusterName = "default"

func main() {

	app := cli.NewApp()
	app.Name = "ecs-watch"
	app.Usage = "keeps track of ecs resources"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{Name: "Patrick Debois",
			Email: "patrick.debois@jedi.be",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "generate",
			Usage: "generates a file",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "cluster", Value: "default", EnvVar: "ECSWATCH_CLUSTER"},
				cli.StringFlag{Name: "region", Value: "eu-west-1", EnvVar: "ECSWATCH_REGION"},
				cli.StringFlag{Name: "template-file", EnvVar: "ECSWATCH_TEMPLATE_FILE"},
				cli.StringFlag{Name: "output-file", EnvVar: "ECSWATCH_OUTPUT_FILE"},
				cli.StringFlag{Name: "notify-container", EnvVar: "ECSWATCH_NOTIFY_CONTAINER"},
				cli.StringFlag{Name: "docker-signal", Value: "SIGHUP", EnvVar: "ECSWATCH_DOCKER_SIGNAL"},
				cli.StringFlag{Name: "docker-container", EnvVar: "ECSWATCH_DOCKER_CONTAINER"},
				cli.StringFlag{Name: "docker-endpoint", Value: "unix:///var/run/docker.sock", EnvVar: "ECSWATCH_DOCKER_ENDPOINT"},
				cli.StringFlag{Name: "watch", Value: "false", EnvVar: "ECSWATCH_WATCH"},
			},
			Action: func(c *cli.Context) error {

				creds := credentials.NewEnvCredentials()

				svc := ecs.New(session.New(), &aws.Config{
					Region:      aws.String(c.String("region")),
					Credentials: creds,
				})
				clusterName = c.String("cluster")
				templateFile := c.String("template-file")
				debug(clusterName)

				options := make(map[string]string)
				options["output-file"] = c.String("output-file")
				options["notify-container"] = c.String("notify-container")

				options["docker-container"] = c.String("docker-container")
				options["docker-endpoint"] = c.String("docker-endpoint")
				options["docker-signal"] = c.String("docker-signal")
				options["watch"] = c.String("watch")

				err := generate(svc, clusterName, templateFile, options)
				if err != nil {
					debug(err.Error())
				}
				return nil
			},
		},
		{
			Name:  "report",
			Usage: "reports all containers and ports",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "cluster", Value: "default", EnvVar: "ECSWATCH_CLUSTER"},
				cli.StringFlag{Name: "region", Value: "eu-west-1", EnvVar: "ECSWATCH_REGION"},
			},
			Action: func(c *cli.Context) error {
				svc := ecs.New(session.New(), &aws.Config{Region: aws.String(c.String("region"))})
				clusterName = c.String("cluster")
				debug(clusterName)
				report(svc, clusterName)
				return nil
			},
		},
	}

	app.Run(os.Args)

}
