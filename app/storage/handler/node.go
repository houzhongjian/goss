package handler

import (
	"net"
	"time"

	"goss.io/goss/lib/logd"
	"goss.io/goss/lib/packet"
	"goss.io/goss/lib/protocol"
)

//connMaster 连接管理节点.
func (this *StorageService) connMaster() {
	//上报节点信息
	conn := this.conn(this.MasterNode)
	pkt := packet.NewNode(packet.NodeTypes_Storage, this.Addr, protocol.NodeAddProtocol)
	_, err := conn.Write(pkt)
	if err != nil {
		this.connMaster()
		return
	}

	for {
		var buf = make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			this.connMaster()
			return
		}
	}
}

func (this *StorageService) conn(node string) net.Conn {
	conn, err := net.Dial("tcp4", node)
	if err != nil {
		logd.Make(logd.Level_WARNING, logd.GetLogpath(), "master节点连接失败，稍后重新连接")
		time.Sleep(time.Second * 1)
		return this.conn(node)
	}

	return conn
}
