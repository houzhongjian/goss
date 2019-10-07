package handler

import (
	"fmt"
	"io"
	"log"
	"net"

	"goss.io/goss/lib"

	"goss.io/goss/lib/protocol"

	"goss.io/goss/lib/packet"

	"goss.io/goss/lib/logd"

	"goss.io/goss/lib/ini"
)

type NodeParams struct {
	Conn  net.Conn
	Types packet.NodeTypes
}

type MasterService struct {
	Conn map[string]NodeParams
	Port string
}

//NewMaster .
func NewMaster() *MasterService {
	return &MasterService{
		Conn: make(map[string]NodeParams),
		Port: fmt.Sprintf(":%d", ini.GetInt("node_port")),
	}
}

//Start.
func (this *MasterService) Start() {
	go NewAdmin()
	this.listen()
	select {}
}

//listen .
func (this *MasterService) listen() {
	listener, err := net.Listen("tcp4", this.Port)
	if err != nil {
		log.Panicln(err)
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			logd.Make(logd.Level_WARNING, logd.GetLogpath(), err.Error())
			return
		}
		ip := conn.RemoteAddr().String()
		logd.Make(logd.Level_INFO, logd.GetLogpath(), "收到来自:"+ip+"的连接请求")
		go this.handler(ip, conn)
	}
}

//handler .
func (this *MasterService) handler(ip string, conn net.Conn) {
	defer conn.Close()
	for {
		pkt, err := packet.ParseNode(conn)
		if err != nil && err == io.EOF {
			logd.Make(logd.Level_WARNING, logd.GetLogpath(), ip+"断开连接")
			//从节点列表中移除.
			RemoveNode(this, ip)
			return
		}

		//判断协议.
		if pkt.Protocol == protocol.NodeAddProtocol {
			//新增节点信息.
			info := Node{
				Types:    pkt.Types,
				IP:       ip,
				SourceIP: pkt.IP,
				CreateAt: lib.Time(),
			}
			NodeInfo = append(NodeInfo, info)

			//记录节点信息.
			this.Conn[pkt.IP] = NodeParams{
				Conn:  conn,
				Types: pkt.Types,
			}

			//新节点上线通知对应的节点.
			if info.Types == packet.NodeTypes_Storage {
				//通知api节点.
				apiList := GetApiList()
				for _, sourceIP := range apiList {
					log.Println("apiSrv:", sourceIP)
					pkt := packet.NewNode(packet.NodeTypes_Storage, pkt.IP, protocol.NodeAddProtocol)
					_, err = this.Conn[sourceIP].Conn.Write(pkt)
					if err != nil {
						logd.Make(logd.Level_WARNING, logd.GetLogpath(), "通知api节点:"+sourceIP+"新增storage节点失败，稍后重新通知")
						RemoveNode(this, ip)
						return
					}
				}
				log.Printf("apiList:%+v\n", apiList)
			}

			if info.Types == packet.NodeTypes_Api {
				//告知新上线的api节点多有的storage节点ip.
				storageList := GetStorageList()

				log.Printf("storageList:%+v\n", storageList)
			}
		}
	}
}
