package handler

import (
	"goss.io/goss/lib/packet"
)

//NodeInfo 节点信息.
var NodeInfo = []Node{}

type Node struct {
	Types    packet.NodeTypes
	IP       string
	CreateAt string
}

//GetStoreList 获取所有的存储节点.
func GetStoreList() []string {
	list := []string{}
	for _, v := range NodeInfo {
		if v.Types == packet.NodeTypes_Storage {
			list = append(list, v.IP)
		}
	}

	return list
}

//获取所有的api节点.
func GetMasterList() []string {
	list := []string{}
	for _, v := range NodeInfo {
		if v.Types == packet.NodeTypes_Api {
			list = append(list, v.IP)
		}
	}

	return list
}
