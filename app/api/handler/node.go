package handler

import (
	"net"
	"time"

	"goss.io/goss/lib/protocol"

	"goss.io/goss/lib/packet"

	"goss.io/goss/lib/logd"
)

//connMaster .
func (this *ApiService) connMaster() {
	conn := this.conn(this.MasterNode)
	pkt := packet.NewNode(packet.NodeTypes_Api, this.Addr, protocol.NodeAddProtocol)
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

//conn .
func (this *ApiService) conn(node string) net.Conn {
	conn, err := net.Dial("tcp4", node)
	if err != nil {
		logd.Make(logd.Level_WARNING, logd.GetLogpath(), "master节点连接失败，稍后重新连接")
		time.Sleep(time.Second * 1)
		return this.conn(node)
	}

	return conn
}
