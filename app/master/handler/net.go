package handler

import (
	"log"
	"net"
	"time"

	"goss.io/goss/app/master/conf"
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
func (this *TcpService) Start() {
	cf := conf.Conf.Node
	for _, addr := range cf.StoreAddrs {
		go this.connStoreNode(addr)
	}
}

//connStoreNode 连接存储节点.
func (this *TcpService) connStoreNode(addr string) {
	for {
		log.Println("开始连接:", addr)
		conn, err := this.connection(addr)
		if err != nil {
			log.Printf("%s:节点连接失败, 尝试重新连接!%+v\n", addr, err)
			time.Sleep(time.Second * 1)
			continue
		}

		this.conn[addr] = conn
		return
	}

	log.Println(addr, "连接成功!")
}

func (this *TcpService) connection(addr string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp4", addr)
	if err != nil {
		return conn, err
	}
	return conn, nil
}

//SelectNode 选择一个存储节点.
func (this *TcpService) SelectNode() net.Conn {
	addr := conf.Conf.Node.StoreAddrs[0]
	return this.conn[addr]
}

//Write tcp 发送消息.
func (this *TcpService) Write(b []byte) (msg []byte, err error) {
	_, err = this.SelectNode().Write(b)
	if err != nil {
		log.Printf("%+v\n", err)
		return msg, err
	}

	for {
		var buf = make([]byte, 1024)
		_, err = this.SelectNode().Read(buf)
		if err != nil {
			log.Printf("%+v\n", err)
			return msg, err
		}

		return buf, nil
	}
}

// var Buf = make(chan []byte, 1024*100)

// const HEADER_LEN = 4
// const HASH_LEN = 32

// type Network struct {
// }

//NewServie .
// func NewServie() {
// 	// srv := &Network{}
// 	// srv.start()
// }

// func (n *Network) start() {
// 	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Conf.Node.Port))
// 	if err != nil {
// 		log.Panicf("%+v\n", err)
// 	}
// 	cf := conf.Conf.Node
// 	msg := fmt.Sprintf("%s(%s:9001) 启动成功!", cf.Name, cf.IP)
// 	println(msg)

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Printf("%+v\n", err)
// 			return
// 		}

// 		log.Println(conn.RemoteAddr().String())
// 		go n.handler(conn)
// 	}
// }

// func (n *Network) handler(conn net.Conn) {
// 	for {
// 		select {
// 		case buf := <-Buf:
// 			msg := fmt.Sprintf("文件已发送给存储节点:%d", len(buf))
// 			h := MD5(buf)
// 			buf = packet(buf, []byte(h))
// 			_, err := conn.Write(buf)
// 			if err != nil {
// 				ip := conn.RemoteAddr().String()
// 				log.Printf("%+v\n", err)
// 				log.Println(ip, "节点断开连接")
// 				return
// 			}

// 			log.Println(msg)

// 			buf = make([]byte, 1024)
// 			_, err = conn.Read(buf)
// 			if err != nil {
// 				log.Printf("%+v\n", err)
// 				return
// 			}
// 			log.Println(string(buf))
// 		}
// 	}
// }

// func packet(content, fileHash []byte) []byte {
// 	buffer := make([]byte, HEADER_LEN+len(content)+len(fileHash))
// 	// 将buffer前面四个字节设置为包长度，大端序
// 	binary.BigEndian.PutUint32(buffer[0:4], uint32(len(content)))
// 	copy(buffer[4:36], fileHash)
// 	copy(buffer[36:], content)
// 	return buffer
// }

// func MD5(body []byte) string {
// 	h := md5.New()
// 	h.Write(body)
// 	b := h.Sum(nil)
// 	return hex.EncodeToString(b)
// }
