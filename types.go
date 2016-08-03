package main

type EcsWatchInfo []EcsWatchInfoItem

type EcsWatchContainerInfo struct {
	InstanceArn   string
	TaskArn       string
	HostPort      int64
	ContainerPort int64
	Image         string
	Name          string
	Status        string
	Environment   map[string]string
	Labels        map[string]string
	Family        string
	Revision      int64
}

type EcsWatchContainersInfo []EcsWatchContainerInfo

type EcsWatchContainerInstanceInfo struct {
	PublicIp    string
	PrivateIp   string
	InstanceArn string
	InstanceId  string
}

type EcsWatchInfoItem struct {
	Cluster       string
	InstanceArn   string
	InstanceId    string
	TaskArn       string
	ContainerPort int64
	HostPort      int64
	Image         string
	Status        string
	Name          string
	PublicIp      string
	PrivateIp     string
	Environment   map[string]string
	Labels        map[string]string
	Family        string
	Revision      int64
}

type EcsWatchEc2InstanceInfo struct {
	PublicIp  string
	PrivateIp string
}

type EcsWatchTasksInfo []EcsWatchTaskInfo

type EcsWatchTaskInfo struct {
	ContainerInstanceArn string
}
