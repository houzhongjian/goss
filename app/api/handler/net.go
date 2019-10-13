package handler

import (
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"goss.io/goss/lib"

	"goss.io/goss/lib/packet"
	"goss.io/goss/lib/protocol"
)

type TcpService struct {
	conn map[string]net.Conn
}

//NewTcpService .
func NewTcpService() *TcpService {
	return &TcpService{
		conn: make(map[string]net.Conn),
	}
}

//Start .
func (this *TcpService) Start(addr string) {
	go this.connStorageNode(addr)
}

//connStorageNode 连接存储节点.
func (this *TcpService) connStorageNode(addr string) {
	for {
		log.Println("开始连接:", addr)
		conn, err := this.connection(addr)
		if err != nil {
			log.Printf("%s:节点连接失败, 尝试重新连接!%+v\n", addr, err)
			time.Sleep(time.Second * 1)
			continue
		}

		this.conn[addr] = conn
		log.Println(addr, "连接成功!")
		return
	}
}

func (this *TcpService) connection(addr string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp4", addr)
	if err != nil {
		return conn, err
	}
	return conn, nil
}

//SelectStoreNode 选择存储节点.
func (this *TcpService) SelectStoreNode() (nodeip string, conn net.Conn) {
	nodeipList := this.SelectNode(1)
	addr := nodeipList[0]
	return addr, this.conn[addr]
}

//SelectNode 选择节点.
//excludeipList 为排除的ip.
func (this *TcpService) SelectNode(nodenum int, excludeipList ...string) []string {
	rand.Seed(time.Now().UnixNano())
	list := []string{}
	for k, _ := range this.conn {
		list = append(list, k)
	}

	nodeipList := []string{}
	num := 0
	for {
		if num >= nodenum {
			break
		}
		index := rand.Int() % len(list)
		addr := list[index]

		//判读当前ip是否需要排除.
		if lib.InArray(addr, excludeipList) {
			continue
		}
		if !lib.InArray(addr, nodeipList) {
			num++
			nodeipList = append(nodeipList, addr)
		}
	}
	return nodeipList
}

//Write tcp 发送消息.
func (this *TcpService) Write(b []byte) (msg []byte, nodeip string, err error) {
	nodeip, conn := this.SelectStoreNode()
	_, err = conn.Write(b)
	if err != nil {
		log.Printf("%+v\n", err)
		return msg, nodeip, err
	}

	for {
		var buf = make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			log.Printf("%+v\n", err)
			return msg, nodeip, err
		}

		return buf, nodeip, nil
	}
}

//Read tcp读取文件.
func (this *TcpService) Read(nodeip, fHash string, bodylen int64) (boby []byte, err error) {
	conn, err := net.Dial("tcp4", nodeip)
	if err != nil {
		log.Printf("%+v\n", err)
		return boby, err
	}

	pkt := packet.New(nil, []byte(fHash), protocol.READ_FILE)
	_, err = conn.Write(pkt)
	if err != nil {
		log.Printf("%+v\n", err)
		return boby, err
	}

	for {
		var buf = make([]byte, bodylen)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			log.Printf("%+v\n", err)
			return boby, err
		}
		return buf, nil
	}
}
