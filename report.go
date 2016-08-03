package main

import "os"
import "fmt"
import "strconv"
import "github.com/olekukonko/tablewriter"

import "github.com/aws/aws-sdk-go/service/ecs"

func report(svc *ecs.ECS, clusterName string) {
	var watchInfo, err = getEcsWatchInfo(svc, clusterName)

	if err != nil {
		debug("[%s] Retrieving report ECS Cluster failed: %s", clusterName, err.Error())
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "HostPort", "ContainerPort", "PublicIp", "PrivateIp", "InstanceId", "Image", "VirtualHost", "Status", "Family", "Revision", "Cluster"})

	for _, infoItem := range *watchInfo {
		fmt.Println(infoItem.Environment)
		table.Append([]string{
			infoItem.Name,
			strconv.FormatInt(infoItem.HostPort, 10),
			strconv.FormatInt(infoItem.ContainerPort, 10),
			infoItem.PublicIp,
			infoItem.PrivateIp,
			infoItem.InstanceId,
			infoItem.Image,
			infoItem.Environment["VIRTUAL_HOST"],
			infoItem.Status,
			infoItem.Family,
			strconv.FormatInt(infoItem.Revision, 10),
			infoItem.Cluster,
		})
	}

	table.Render()

}
