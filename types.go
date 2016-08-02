package main

type EcsWatchInfo []EcsWatchInfoItem

type EcsWatchContainerInfo struct {
	instanceArn   string
	taskArn       string
	hostPort      int64
	containerPort int64
	imageName     string
	name          string
}

type EcsWatchContainersInfo []EcsWatchContainerInfo

type EcsWatchContainerInstanceInfo struct {
	publicIp    string
	privateIp   string
	instanceArn string
	instanceId  string
}

type EcsWatchInfoItem struct {
	instanceArn   string
	taskArn       string
	containerPort int64
	hostPort      int64
	imageName     string
	name          string
	publicIp      string
	privateIp     string
}

type EcsWatchEc2InstanceInfo struct {
	publicIp  string
	privateIp string
}

type EcsWatchTasksInfo []EcsWatchTaskInfo

type EcsWatchTaskInfo struct {
	containerInstanceArn string
}
