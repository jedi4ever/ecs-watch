package main

//import "os"
//import "strconv"
import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecs"
	"time"
)

func track(svc *ecs.ECS, clusterName string, options EcsWatchTrackOptions) error {

	lastKnownInfo, err := getEcsWatchInfo(svc, clusterName)

	if err != nil {
		debug("[%s] Retrieving report ECS Cluster failed: %s", clusterName, err.Error())
		return err
	}

	tickChan := time.NewTicker(options.TrackInterval).C

	doneChan := make(chan bool)

	for {
		select {
		case <-tickChan:
			debug("tick")

			currentInfo, err := getEcsWatchInfo(svc, clusterName)

			if err != nil {
				debug("[%s] Retrieving report ECS Cluster failed: %s", clusterName, err.Error())
				return err
			}

			if hasInfoChanged(*lastKnownInfo, *currentInfo) {

				if options.TemplateGenerate {
					err := templateGenerate(*currentInfo, options)
					if err != nil {
						debug("[%s] Generating template %s failed : %s", clusterName, options.TemplateInputFile, err.Error())
						return err
					}
				}

				if options.DockerNotify {
					err := dockerSignal(options.DockerSignal, options.DockerContainer, options.DockerEndpoint)

					if err != nil {
						debug("[%s] Error sending docker signal %s failed : %s", clusterName, options.DockerContainer, err.Error())
						return err
					}
				}
			}
			lastKnownInfo = currentInfo

		case <-doneChan:
			fmt.Println("Done")
			return nil
		}
	}

}

func hasInfoChanged(lastKnownInfo EcsWatchInfo, currentKnownInfo EcsWatchInfo) (changed bool) {
	return true
}
