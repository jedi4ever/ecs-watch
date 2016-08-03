package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Retrieves the complete information about an ecs cluster
func getEcsWatchInfo(svc *ecs.ECS, clusterName string) (ecsWatchInfo *EcsWatchInfo, err error) {

	ecsWatchInfo2 := EcsWatchInfo{}
	debug("[%s] Retrieving Container Info about ECS cluster ", clusterName)
	containersInfo, err := getEcsWatchContainersInfo(svc, clusterName)

	if err != nil {
		debug("[%s] Retrieving Task Info for ECS Cluster failed: %s", clusterName, err.Error())
		return nil, err
	}

	// For each task , retrieve the container info
	for _, containerInfo := range *containersInfo {
		ecsWatchContainerInstanceInfo, err := getEcsWatchContainerInstanceInfo(svc, containerInfo.InstanceArn)
		if err != nil {
			debug("[%s] Retrieving ContainerInstance Info for ECS Cluster failed: %s", clusterName, err.Error())
			return nil, err
		}

		ecsWatchInfoItem := &EcsWatchInfoItem{
			Cluster:       clusterName,
			PublicIp:      ecsWatchContainerInstanceInfo.PublicIp,
			PrivateIp:     ecsWatchContainerInstanceInfo.PrivateIp,
			InstanceArn:   containerInfo.InstanceArn,
			InstanceId:    ecsWatchContainerInstanceInfo.InstanceId,
			TaskArn:       containerInfo.TaskArn,
			Status:        containerInfo.Status,
			Image:         containerInfo.Image,
			Name:          containerInfo.Name,
			HostPort:      containerInfo.HostPort,
			ContainerPort: containerInfo.ContainerPort,
		}

		ecsWatchInfo2 = append(ecsWatchInfo2, *ecsWatchInfoItem)

	}

	return &ecsWatchInfo2, nil
}

func getEcsWatchContainersInfo(svc *ecs.ECS, clusterName string) (ecsWatchContainersInfo *EcsWatchContainersInfo, err error) {

	// Lets first get all tasks
	params := &ecs.ListTasksInput{
		Cluster: aws.String(clusterName),
		//ContainerInstance: aws.String("String"),
		//DesiredStatus:     aws.String("DesiredStatus"),
		//Family:            aws.String("String"),
		//MaxResults:        aws.Int64(1),
		//NextToken:         aws.String("String"),
		//ServiceName:       aws.String("String"),
		//StartedBy:         aws.String("String"),
	}
	resp, err := svc.ListTasks(params)

	if err != nil {
		debug(err.Error())
		return
	}

	ecsWatchContainersInfo2 := EcsWatchContainersInfo{}

	// For each taskArn
	for _, taskArn := range resp.TaskArns {
		task, err := describeTask(svc, clusterName, *taskArn)
		if err != nil {
			debug(err.Error())
			return nil, err
		}

		// Now for each container in this task
		for _, container := range task.Containers {

			// Create a new ContainerInfo
			ecsWatchContainerInfo := EcsWatchContainerInfo{}
			ecsWatchContainerInfo.Environment = make(map[string]string)
			ecsWatchContainerInfo.Labels = make(map[string]string)
			ecsWatchContainerInfo.Name = *container.Name
			ecsWatchContainerInfo.Status = *container.LastStatus

			// Find network Binding
			for _, binding := range container.NetworkBindings {
				ecsWatchContainerInfo.ContainerPort = *binding.ContainerPort
				ecsWatchContainerInfo.HostPort = *binding.HostPort
				debug("%v|%v|%s", *binding.ContainerPort, *binding.HostPort, *binding.Protocol)
			}

			ecsWatchContainerInfo.InstanceArn = *task.ContainerInstanceArn
			debug("instanceArn = %s", ecsWatchContainerInfo.InstanceArn)

			// Find task definition of this container
			taskDefinition, err := describeTaskDefinition(svc, task)

			if err != nil {
				debug(err.Error())
				return nil, err
			}

			for _, containerDefinition := range taskDefinition.ContainerDefinitions {
				if *containerDefinition.Name == *container.Name {
					for _, pair := range containerDefinition.Environment {
						debug("found Environment var %s = %s", *pair.Name, *pair.Value)
						ecsWatchContainerInfo.Environment[*pair.Name] = *pair.Value
					}
					ecsWatchContainerInfo.Image = *containerDefinition.Image

					// Iterate maps
					for k, v := range containerDefinition.DockerLabels {
						debug("found Label var %s = %s", k, *v)
						ecsWatchContainerInfo.Labels[k] = *v //debug(containerDefinition.Environment)
					}
					//fmt.Println(container)
				}

			}

			ecsWatchContainersInfo2 = append(ecsWatchContainersInfo2, ecsWatchContainerInfo)

		}

	}

	return &ecsWatchContainersInfo2, nil

}

func describeTaskDefinition(svc *ecs.ECS, task *ecs.Task) (taskDefinition *ecs.TaskDefinition, err error) {

	params := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(*task.TaskDefinitionArn), // Required
	}

	resp, err := svc.DescribeTaskDefinition(params)

	if err != nil {
		debug(err.Error())
		return nil, err
	}
	return resp.TaskDefinition, nil

}

func getEcsWatchEc2InstanceInfo(svcEc2 *ec2.EC2, instanceID string) (ec2Info *EcsWatchEc2InstanceInfo, err error) {

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(instanceID),
				},
			},
		},
	}

	resp, err := svcEc2.DescribeInstances(params)

	if err != nil {
		return nil, err
	}

	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			//if inst.PrivateIpAddress != nil {
			return &EcsWatchEc2InstanceInfo{
				PrivateIp: *inst.PrivateIpAddress,
				PublicIp:  *inst.PublicIpAddress,
			}, nil
			//}
		}
	}

	return nil, errors.New("No vm found with that maches instanceID " + instanceID)
}

// Retrieve information of a containerInstance
func getEcsWatchContainerInstanceInfo(svc *ecs.ECS, instanceArn string) (containerInstanceInfo *EcsWatchContainerInstanceInfo, err error) {

	params := &ecs.DescribeContainerInstancesInput{
		ContainerInstances: []*string{ // Required
			aws.String(instanceArn),
		},
		Cluster: aws.String(clusterName),
	}
	resp, err := svc.DescribeContainerInstances(params)

	if err != nil {
		debug(err.Error())
		return nil, err
	}

	// Get the instance Id
	instanceId := *resp.ContainerInstances[0].Ec2InstanceId

	svcEc2 := ec2.New(session.New(), &aws.Config{Region: aws.String("eu-west-1")})

	ec2Info, err := getEcsWatchEc2InstanceInfo(svcEc2, (*resp.ContainerInstances[0].Ec2InstanceId))

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	containerInstanceInfo = &EcsWatchContainerInstanceInfo{
		InstanceArn: instanceArn,
		InstanceId:  instanceId,
		PrivateIp:   ec2Info.PrivateIp,
		PublicIp:    ec2Info.PublicIp,
	}

	return containerInstanceInfo, nil

}

// Get the info about a single Ecs task
func describeTask(svc *ecs.ECS, clusterName string, taskArn string) (task *ecs.Task, err error) {

	params := &ecs.DescribeTasksInput{
		Tasks: []*string{ // Required
			aws.String(taskArn), // Required
			// More values...
		},
		Cluster: aws.String(clusterName),
	}

	resp, err := svc.DescribeTasks(params)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp.Tasks[0], nil

}
