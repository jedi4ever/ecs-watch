package main

type EcsWatchInfo []EcsWatchInfoItem

type EcsWatchContainerInfo struct {
	InstanceArn   string
	TaskArn       string
	HostPort      int64
	ContainerPort int64
	ImageName     string
	Name          string
}

type EcsWatchContainersInfo []EcsWatchContainerInfo

type EcsWatchContainerInstanceInfo struct {
	PublicIp    string
	PrivateIp   string
	InstanceArn string
	InstanceId  string
}

type EcsWatchInfoItem struct {
	InstanceArn   string
	TaskArn       string
	ContainerPort int64
	HostPort      int64
	ImageName     string
	Name          string
	PublicIp      string
	PrivateIp     string
}

type EcsWatchEc2InstanceInfo struct {
	PublicIp  string
	PrivateIp string
}

type EcsWatchTasksInfo []EcsWatchTaskInfo

type EcsWatchTaskInfo struct {
	ContainerInstanceArn string
}
