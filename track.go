package main

//import "os"
//import "strconv"
import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecs"
	"reflect"
	"time"
)

func track(svc *ecs.ECS, clusterName string, options EcsWatchTrackOptions) error {

	lastKnownInfo, err := getEcsWatchInfo(svc, clusterName)

	if err != nil {
		debug("[%s] Retrieving report ECS Cluster failed: %s", clusterName, err.Error())
		return err
	}

	err = handleOnce(*lastKnownInfo, clusterName, options)
	if err != nil {
		return err
	}

	if options.OnlyOnce {
		return nil
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
				err := handleOnce(*currentInfo, clusterName, options)
				if err != nil {
					return err
				}
				debug("info has changed")
				lastKnownInfo = currentInfo
			}

		case <-doneChan:
			fmt.Println("Done")
			return nil
		}
	}

}

func handleOnce(currentInfo EcsWatchInfo, clusterName string, options EcsWatchTrackOptions) error {
	if options.TemplateGenerate {
		debug("generating template")
		err := templateGenerate(currentInfo, options)
		if err != nil {
			debug("[%s] Generating template %s failed : %s", clusterName, options.TemplateInputFile, err.Error())
			return err
		}
	}

	if options.DockerNotify {
		debug("notifying docker")
		err := dockerSignal(options.DockerSignal, options.DockerContainer, options.DockerEndpoint)

		if err != nil {
			debug("[%s] Error sending docker signal %s failed : %s", clusterName, options.DockerContainer, err.Error())
			return err
		}
	}
	return nil
}

func hasInfoChanged(lastKnownInfo EcsWatchInfo, currentKnownInfo EcsWatchInfo) (changed bool) {
	// http://stackoverflow.com/questions/24534072/how-to-compare-struct-slice-map-are-equal
	// https://www.reddit.com/r/golang/comments/369o32/map_with_custom_equality/
	return !reflect.DeepEqual(lastKnownInfo, currentKnownInfo)
}
