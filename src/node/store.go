package node

import (
	"fmt"
	"log"
	"net"
	"time"

	"pandaschool.net/goss/src/conf"
)

//Store 存储节点.
type Store struct {
}

//NewStore .
func NewStore() *Store {
	store := &Store{}
	return store
}

func (s *Store) Start() {
	log.Println("存储节点创建成功!")
	log.Println("准备上包节点信息给主节点!")
	s.report()

	s.listen()
}

//report 上报节点信息到主节点.
func (s *Store) report() {

}

func (s *Store) listen() {
	listen, err := net.Listen("tcp", conf.Conf.Node.Store[0])
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}
	for {
		conn, err := listen.Accept()
		ip := conn.RemoteAddr().String()
		if err != nil {
			log.Printf("%+v\n", err)
			log.Println(ip, ":断开连接")
			return
		}
		log.Println(ip, "主节点已连接")
		go s.handler(conn)
	}
}

func (s *Store) handler(conn net.Conn) {
	defer conn.Close()
	for {
		node := conf.Conf.Node.Name
		msg := fmt.Sprintf("我是%s, 我当前的时间为:%s", node, time.Now().Format("2006-01-02 15:04:05"))
		ip := conn.RemoteAddr().String()
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Printf("%+v\n", err)
			log.Println(ip, ":主节点断开连接，等待主节点重新连接")
			return
		}
		time.Sleep(time.Second * 1)
	}
}
