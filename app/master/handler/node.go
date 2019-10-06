package handler

import (
	"log"

	"goss.io/goss/lib/packet"
)

//NodeInfo 节点信息.
var NodeInfo = []Node{}

type Node struct {
	Types    packet.NodeTypes
	IP       string
	SourceIP string //所属ip.
	CreateAt string
}

//GetStoreList 获取所有的存储节点.
func GetStorageList() []string {
	list := []string{}
	for _, v := range NodeInfo {
		if v.Types == packet.NodeTypes_Storage {
			list = append(list, v.SourceIP)
		}
	}

	return list
}

//获取所有的api节点.
func GetApiList() []string {
	list := []string{}
	for _, v := range NodeInfo {
		if v.Types == packet.NodeTypes_Api {
			list = append(list, v.SourceIP)
		}
	}

	return list
}

//RemoveNode 移除某一个切片.
func RemoveNode(ip string) {
	log.Println("需要删除的ip:", ip)
	for index, v := range NodeInfo {
		if v.IP == ip {
			NodeInfo = append(NodeInfo[:index], NodeInfo[index+1:]...)
		}
	}
	log.Printf("NodeInfo:%+v\n", NodeInfo)
}
