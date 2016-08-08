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
	"time"
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
			Name:  "track",
			Usage: "tracks changes inside an ecs cluster",
			Flags: []cli.Flag{
				cli.DurationFlag{Name: "track-interval", Value: time.Second * 5, EnvVar: "ECSWATCH_TRACK_INTERVAL"},
				cli.BoolFlag{Name: "only-once", EnvVar: "ECSWATCH_ONLY_ONCE"},

				// Future filter items
				cli.StringFlag{Name: "ecs-cluster", Value: "default", EnvVar: "ECSWATCH_ECS_CLUSTER, ECS_CLUSTER"},
				cli.StringFlag{Name: "aws-region", Value: "eu-west-1", EnvVar: "ECSWATCH_AWS_REGION, AWS_REGION"},

				/*
					cli.StringFlag{Name: "log-format", Value: "json", EnvVar: "ECSWATCH_LOG_FORMAT"},
					cli.BoolFlag{Name: "debug", EnvVar: "ECSWATCH_DEBUG"},
					cli.BoolFlag{Name: "silent", EnvVar: "ECSWATCH_SILENT"},
				*/

				/*
					cli.BoolFlag{Name: "route53-update", EnvVar: "ECSWATCH_UPDATE_ROUTE53"},
					cli.StringFlag{Name: "route53-zone", EnvVar: "ECSWATCH_REGISTER_ROUTE53_ZONE"},
					cli.StringFlag{Name: "route53-template", EnvVar: "ECSWATCH_REGISTER_ROUTE53_TEMPLATE"},
				*/

				cli.BoolFlag{Name: "template-generate", EnvVar: "ECSWATCH_TEMPLATE_FILE"},
				cli.StringFlag{Name: "template-input-file", EnvVar: "ECSWATCH_TEMPLATE_INPUT_FILE"},
				cli.StringFlag{Name: "template-output-file", EnvVar: "ECSWATCH_TEMPLATE_OUTPUT_FILE"},

				cli.BoolFlag{Name: "docker-notify", EnvVar: "ECSWATCH_NOTIFY_CONTAINER"},
				cli.StringFlag{Name: "docker-signal", Value: "SIGHUP", EnvVar: "ECSWATCH_DOCKER_SIGNAL"},
				cli.StringFlag{Name: "docker-container", EnvVar: "ECSWATCH_DOCKER_CONTAINER"},
				cli.StringFlag{Name: "docker-endpoint", Value: "unix:///var/run/docker.sock", EnvVar: "ECSWATCH_DOCKER_ENDPOINT"},
				//cli.StringSliceFlag{Name: "filter", Value: &cli.StringSlice{"Name=*"}, EnvVar: "ECSWATCH_FILTER"},

				// Future - datadog metrics , events
				// Future - exec , ip lookup
				// Future - cloudlogs
				// Future - sns endpoint
			},
			Action: func(c *cli.Context) error {

				creds := credentials.NewEnvCredentials()

				svc := ecs.New(session.New(), &aws.Config{
					Region:      aws.String(c.String("aws-region")),
					Credentials: creds,
				})

				options := EcsWatchTrackOptions{}

				clusterName = c.String("ecs-cluster")
				options.EcsCluster = c.String("ecs-cluster")
				options.AwsRegion = c.String("aws-region")

				options.TrackInterval = c.Duration("track-interval")
				options.OnlyOnce = c.Bool("only-once")

				options.TemplateGenerate = c.Bool("template-generate")
				options.TemplateOutputFile = c.String("template-output-file")
				options.TemplateInputFile = c.String("template-input-file")

				options.DockerNotify = c.Bool("docker-notify")
				options.DockerContainer = c.String("docker-container")
				options.DockerEndpoint = c.String("docker-endpoint")
				options.DockerSignal = c.String("docker-signal")

				err := track(svc, clusterName, options)

				if err != nil {
					cli.NewExitError(err.Error(), -1)
					debug(err.Error())
				}
				return nil
			},
		},
		{
			Name:  "report",
			Usage: "reports all containers and ports",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "ecs-cluster", Value: "default", EnvVar: "ECSWATCH_ECS_CLUSTER"},
				cli.StringFlag{Name: "aws-region", Value: "eu-west-1", EnvVar: "ECSWATCH_AWS_REGION"},
			},
			Action: func(c *cli.Context) error {
				svc := ecs.New(session.New(), &aws.Config{Region: aws.String(c.String("aws-region"))})
				clusterName = c.String("ecs-cluster")
				debug(clusterName)
				err := report(svc, clusterName)
				if err != nil {
					cli.NewExitError(err.Error(), -1)
					debug(err.Error())
				}
				return nil
			},
		},
	}

	app.Run(os.Args)

}

type EcsWatchTrackOptions struct {
	EcsCluster         string
	AwsRegion          string
	TrackInterval      time.Duration
	OnlyOnce           bool
	TemplateGenerate   bool
	TemplateOutputFile string
	TemplateInputFile  string
	DockerNotify       bool
	DockerContainer    string
	DockerEndpoint     string
	DockerSignal       string
}
