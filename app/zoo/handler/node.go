package handler

import (
	"time"

	"goss.io/goss/lib/packet"
)

//NodeInfo 节点信息.
var NodeInfo = []Node{}

type Node struct {
	Types    packet.NodeTypes
	IP       string
	CreateAt time.Time
}

//GetStoreList 获取所有的存储节点.
func GetStoreList() []string {
	list := []string{}
	for _, v := range NodeInfo {
		if v.Types == packet.NodeTypes_Store {
			list = append(list, v.IP)
		}
	}

	return list
}

//获取所有的master节点.
func GetMasterList() []string {
	list := []string{}
	for _, v := range NodeInfo {
		if v.Types == packet.NodeTypes_Master {
			list = append(list, v.IP)
		}
	}

	return list
}
