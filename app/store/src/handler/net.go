package handler

import (
	"fmt"
	"log"
	"net"

	"goss.io/goss/app/store/src/conf"
)

type StoreService struct {
	Port string
}

func NewStoreService() *StoreService {
	s := &StoreService{
		Port: fmt.Sprintf(":%d", conf.Conf.Node.Port),
	}
	return s
}

//Start .
func (this *StoreService) Start() {
	this.listen()
}

//listen .
func (this *StoreService) listen() {
	listener, err := net.Listen("tcp4", this.Port)
	if err != nil {
		log.Printf("%s 端口监听失败!%+v\n", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("%+v\n", err)
			return
		}

		go this.handler(conn)
	}
}

func (this *StoreService) handler(conn net.Conn) {
	defer conn.Close()
	for {
		var buf = make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			log.Printf("%+v\n", err)
			return
		}

		log.Println(string(buf))

		conn.Write(buf)
	}
}

// const HEADER_LEN = 4
// const HASH_LEN = 32

// //NewHandler .
// func NewHandler() {

// }

// func conn(ip string) (net.Conn, error) {
// 	conn, err := net.Dial("tcp", ip)
// 	if err != nil {
// 		log.Printf("%+v\n", err)
// 		return conn, err
// 	}

// 	log.Println(ip, "连接成功")
// 	return conn, nil
// }

// func connMaster(ip string) {
// 	log.Println(ip)
// 	conn, err := conn(ip)
// 	if err != nil {
// 		log.Printf("连接失败，准备尝试重新连接:%+v\n", err)
// 		connMaster(ip)
// 		return
// 	}
// 	for {
// 		var header = make([]byte, HEADER_LEN)
// 		_, err = io.ReadFull(conn, header)
// 		if err != nil {
// 			log.Printf("%+v\n", err)
// 			return
// 		}
// 		length := int(binary.BigEndian.Uint32(header))
// 		log.Println("收到来自主节点的数据请求，数据大小为：", length)

// 		var fhash = make([]byte, HASH_LEN)
// 		_, err = io.ReadFull(conn, fhash)
// 		if err != nil {
// 			log.Printf("%+v\n", err)
// 			return
// 		}
// 		fileHash := string(fhash)

// 		//读取文件.
// 		var buf = make([]byte, length)
// 		_, err = io.ReadFull(conn, buf)
// 		if err != nil {
// 			log.Printf("%+v\n", err)
// 			return
// 		}

// 		if fileHash != MD5(buf) {
// 			_, err = conn.Write([]byte("文件hash不一致"))
// 			if err != nil {
// 				log.Printf("%+v\n", err)
// 				return
// 			}
// 			return
// 		}

// 		fname := fmt.Sprintf("./tmp/%s.jpg", uuid.New().String())
// 		f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
// 		if err != nil {
// 			f.Close()
// 			log.Printf("%+v\n", err)
// 			return
// 		}
// 		_, err = f.Write(buf)
// 		if err != nil {
// 			log.Printf("%+v\n", err)
// 			return
// 		}
// 		f.Close()

// 		msg := fname + " 创建成功!"
// 		log.Println(msg)
// 		_, err = conn.Write([]byte(msg))
// 		if err != nil {
// 			log.Printf("%+v\n", err)
// 			return
// 		}
// 	}
// }

// func MD5(body []byte) string {
// 	h := md5.New()
// 	h.Write(body)
// 	b := h.Sum(nil)
// 	return hex.EncodeToString(b)
// }
