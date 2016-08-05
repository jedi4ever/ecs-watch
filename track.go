package main

//import "os"
//import "strconv"
import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecs"
	"time"
)

func track(svc *ecs.ECS, clusterName string, options EcsWatchTrackOptions) error {

	var watchInfo, err = getEcsWatchInfo(svc, clusterName)

	if err != nil {
		debug("[%s] Retrieving report ECS Cluster failed: %s", clusterName, err.Error())
		return err
	}

	tickChan := time.NewTicker(options.TrackInterval).C

	doneChan := make(chan bool)

	result, err := templateGenerate(*watchInfo, options)

	if err != nil {
		debug("[%s] Generating template %s failed : %s", clusterName, options.TemplateInputFile, err.Error())
		return err
	}
	prevResult := result

	for {
		select {
		case <-tickChan:
			debug("tick")
			result, err := templateGenerate(*watchInfo, options)
			if err != nil {
				debug("[%s] Generating template %s failed : %s", clusterName, options.TemplateInputFile, err.Error())
				return err
			}

			dockerSignal(options.DockerSignal, options.DockerContainer, options.DockerEndpoint)

			if result != prevResult {
				prevResult = result
				debug("******** CHANGED DETECTED ****")
				dockerSignal(options.DockerSignal, options.DockerContainer, options.DockerEndpoint)
			}

		case <-doneChan:
			fmt.Println("Done")
			return nil
		}
	}

}
