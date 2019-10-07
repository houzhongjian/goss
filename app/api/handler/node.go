package handler

import (
	"io"
	"net"
	"time"

	"goss.io/goss/lib/protocol"

	"goss.io/goss/lib/packet"

	"goss.io/goss/lib/logd"
)

//connMaster .
func (this *ApiService) connMaster() {
	conn := this.conn(this.MasterNode)
	logd.Make(logd.Level_INFO, logd.GetLogpath(), "master节点连接成功，准备开始上报节点信息")
	pkt := packet.NewNode(packet.NodeTypes_Api, this.Addr, protocol.NodeAddProtocol)
	_, err := conn.Write(pkt)
	if err != nil {
		this.connMaster()
		return
	}
	logd.Make(logd.Level_INFO, logd.GetLogpath(), "上报节点信息成功")

	for {
		pkt, err := packet.ParseNode(conn)
		if err != nil && err == io.EOF {
			logd.Make(logd.Level_WARNING, logd.GetLogpath(), "master节点断开连接，稍后重新连接master节点")
			this.connMaster()
			return
		}

		//判断协议类型.
		if pkt.Protocol == protocol.NodeAddProtocol {
			//新增节点.
			logd.Make(logd.Level_INFO, logd.GetLogpath(), "新增存储节点:"+pkt.IP)
			this.Tcp.Start(pkt.IP)
		}

		if pkt.Protocol == protocol.NodelDelProtocol {
			//删除节点.
			logd.Make(logd.Level_INFO, logd.GetLogpath(), "收到master节点要求与:"+pkt.IP+"节点断开的消息")
			if err := this.Tcp.conn[pkt.IP].Close(); err != nil {
				logd.Make(logd.Level_INFO, logd.GetLogpath(), "断开与:"+pkt.IP+"节点的连接失败")
				return
			}
			logd.Make(logd.Level_INFO, logd.GetLogpath(), "断开成功")
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
