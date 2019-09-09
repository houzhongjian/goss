package handler

import (
	"log"
	"net"

	"pandaschool.net/goss/src/conf"
)

//Gateway 网关.
type Gateway struct {
	StoreNode []string
}

func NewMaster() *Gateway {
	gwy := &Gateway{
		StoreNode: conf.Conf.Node.Store,
	}
	return gwy
}

//Start .
func (g *Gateway) Start() {
	log.Println("开始读取存储节点信息")
	for _, addr := range g.StoreNode {
		go g.ConnStoreNode(addr)
	}
	go g.dashboard()
	select {}
}

func (g *Gateway) ConnStoreNode(addr string) {
	log.Println("准备开始连接:", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(addr, "连接失败")
		log.Println("系统尝试重新连接：", addr)
		//todo 将断开的存储节点从可用节点的列表中删除
		g.ConnStoreNode(addr)
		return
	}
	for {
		buffer := make([]byte, 1024)
		_, err = conn.Read(buffer)
		if err != nil {
			//todo 将断开的存储节点从可用节点的列表中删除
			log.Println("存储节点:", addr, "下线")
			log.Println("系统尝试重新连接：", addr)
			g.ConnStoreNode(addr)
			return
		}

		log.Println(string(buffer))
	}
}

//handle .
func (g *Gateway) handle(conn net.Conn) {
	defer conn.Close()
	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			log.Printf("%+v\n", err)
			return
		}

		log.Println(string(buffer))
	}
}

//存储服务控制台.
func (g *Gateway) dashboard() {

}
