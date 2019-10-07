package handler

import (
	"goss.io/goss/lib/packet"
	"goss.io/goss/lib/protocol"
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
func RemoveNode(this *MasterService, ip string) {
	//根据访问ip获取节点ip.
	for index, v := range NodeInfo {
		if v.IP == ip {
			//通知对应的节点与故障节点断开连接.
			apiList := GetApiList()
			for _, apinodeIp := range apiList {
				pkt := packet.NewNode(packet.NodeTypes_Api, v.SourceIP, protocol.NodelDelProtocol)
				this.Conn[apinodeIp].Conn.Write(pkt)
			}

			//从NodeInfo中移除当前.
			NodeInfo = append(NodeInfo[:index], NodeInfo[index+1:]...)

			//删除对应的连接数据.
			delete(this.Conn, v.SourceIP)
		}
	}
}
