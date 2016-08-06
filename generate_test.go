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
		if len(group["www"]) != 1 {
			t.Error("Expected only one item in the www group")
		}
	} else {
		t.Error("Expected www group to be available")
	}

}

func TestGroupByVirtualHost_twosamevirtualhostsenvvar(t *testing.T) {

	ecsWatchInfoItem := &EcsWatchInfoItem{
		Environment: map[string]string{
			"VIRTUAL_HOST": "www",
		},
	}

	var ecsWatchInfo EcsWatchInfo
	ecsWatchInfo = append(ecsWatchInfo, *ecsWatchInfoItem)
	ecsWatchInfo = append(ecsWatchInfo, *ecsWatchInfoItem)

	group := groupByVirtualHost(ecsWatchInfo)

	if len(group) != 1 {
		t.Error("Expected 1 group got ", len(group))
	}

	if _, ok := group["www"]; ok {
		if len(group["www"]) != 2 {
			t.Error("Expected two item2 in the www group")
		}
	} else {
		t.Error("Expected www group to be available")
	}

}

func TestGroupByVirtualHost_twodifferentvirtualhostsenvvar(t *testing.T) {

	ecsWatchInfoItem1 := &EcsWatchInfoItem{
		Environment: map[string]string{
			"VIRTUAL_HOST": "www1",
		},
	}

	ecsWatchInfoItem2 := &EcsWatchInfoItem{
		Environment: map[string]string{
			"VIRTUAL_HOST": "www2",
		},
	}

	var ecsWatchInfo EcsWatchInfo
	ecsWatchInfo = append(ecsWatchInfo, *ecsWatchInfoItem1)
	ecsWatchInfo = append(ecsWatchInfo, *ecsWatchInfoItem2)

	group := groupByVirtualHost(ecsWatchInfo)

	if len(group) != 2 {
		t.Error("Expected 1 group got ", len(group))
	}

	if _, ok := group["www1"]; ok {
		if len(group["www1"]) != 1 {
			t.Error("Expected one item in the www1 group")
		}
	} else {
		t.Error("Expected www1 group to be available")
	}

}
