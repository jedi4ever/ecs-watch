package main

import "testing"

func TestGroupByVirtualHost_empty(t *testing.T) {

	//ecsWatchInfoItem := &EcsWatchInfoItem{}

	var ecsWatchInfo EcsWatchInfo
	//ecsWatchInfo = append(ecsWatchInfo, *ecsWatchInfoItem)

	group := groupByVirtualHost(ecsWatchInfo)

	if len(group) != 0 {
		t.Error("Expected nil group got ", group)
	}

}

func TestGroupByVirtualHost_noenvironmentvars(t *testing.T) {

	ecsWatchInfoItem := &EcsWatchInfoItem{
		Name: "bla",
	}

	var ecsWatchInfo EcsWatchInfo
	ecsWatchInfo = append(ecsWatchInfo, *ecsWatchInfoItem)

	group := groupByVirtualHost(ecsWatchInfo)

	if len(group) != 0 {
		t.Error("Expected nil group got ", group)
	}

}

func TestGroupByVirtualHost_novirtualhostsenvvar(t *testing.T) {

	ecsWatchInfoItem := &EcsWatchInfoItem{
		Environment: map[string]string{
			"var1": "value1",
		},
	}

	var ecsWatchInfo EcsWatchInfo
	ecsWatchInfo = append(ecsWatchInfo, *ecsWatchInfoItem)

	group := groupByVirtualHost(ecsWatchInfo)

	if len(group) != 0 {
		t.Error("Expected nil group got ", group)
	}

}

func TestGroupByVirtualHost_virtualhostsenvvar(t *testing.T) {

	ecsWatchInfoItem := &EcsWatchInfoItem{
		Environment: map[string]string{
			"VIRTUAL_HOST": "www",
		},
	}

	var ecsWatchInfo EcsWatchInfo
	ecsWatchInfo = append(ecsWatchInfo, *ecsWatchInfoItem)

	group := groupByVirtualHost(ecsWatchInfo)

	if len(group) != 1 {
		t.Error("Expected 1 group got ", len(group))
	}

	if _, ok := group["www"]; ok {
		//do something here
	} else {
		t.Error("Expected www group to be available")
	}

}
