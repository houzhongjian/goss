package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"goss.io/goss/lib/logd"

	"goss.io/goss/lib/ini"

	"goss.io/goss/lib/protocol"

	"goss.io/goss/app/storage/conf"
	"goss.io/goss/lib"
	"goss.io/goss/lib/packet"
)

type StorageService struct {
	Port       string
	Addr       string
	MasterNode string
}

func NewStorageService() *StorageService {
	s := &StorageService{
		Port:       fmt.Sprintf(":%d", conf.Conf.Node.Port),
		Addr:       fmt.Sprintf("%s:%d", ini.GetString("node_ip"), ini.GetInt("node_port")),
		MasterNode: ini.GetString("master_node"),
	}
	return s
}

//Start .
func (this *StorageService) Start() {
	this.checkStoragePath()
	go this.connMaster()
	this.listen()
}

//connMaster 连接管理节点.
func (this *StorageService) connMaster() {
	//上报节点信息
	conn := this.conn(this.MasterNode)
	pkt := packet.NewNode(packet.NodeTypes_Storage, this.Addr, protocol.NodeAddProtocol)
	conn.Write(pkt)
}

func (this *StorageService) conn(node string) net.Conn {
	conn, err := net.Dial("tcp4", node)
	if err != nil {
		logd.Make(logd.Level_WARNING, logd.GetLogpath(), "master节点连接失败，稍后重新连接")
		time.Sleep(time.Second * 1)
		this.conn(node)
	}

	return conn
}

//checkStoragePath 检查存储路径.
func (this *StorageService) checkStoragePath() {
	if !lib.IsExists(conf.Conf.Node.StorageRoot) {
		//创建存储文件夹.
		if err := os.Mkdir(conf.Conf.Node.StorageRoot, 0777); err != nil {
			log.Panicf("%+v\n", err)
		}
	}
}

//listen .
func (this *StorageService) listen() {
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

func (this *StorageService) handler(conn net.Conn) {
	defer conn.Close()
	for {
		pkt, err := packet.Parse(conn)
		if err != nil {
			log.Printf("%+v\n", err)
			return
		}

		//判断协议号.
		if pkt.Protocol == protocol.WriteFileProrocol {
			//计算文件hash.
			fHash := lib.FileHash(pkt.Body)
			//验证文件是否损坏.
			log.Println("计算的hash为:", fHash)
			if fHash != pkt.Hash {
				log.Println("文件不一致")
				conn.Write([]byte("fail"))
				return
			}

			fPath := conf.Conf.Node.StorageRoot + fHash
			err = ioutil.WriteFile(fPath, pkt.Body, 0777)
			if err != nil {
				log.Printf("%+v\n", err)
				return
			}
			conn.Write([]byte(fHash))
		}

		if pkt.Protocol == protocol.ReadFileProrocol {
			log.Println("读取文件：", pkt.Hash)
			//读取文件.
			fpath := conf.Conf.Node.StorageRoot + pkt.Hash
			b, err := ioutil.ReadFile(fpath)
			if err != nil {
				log.Printf("%+v\n", err)
				return
			}

			log.Println("文件大小为:", len(b))

			_, err = conn.Write(b)
			if err != nil {
				log.Printf("%+v\n", err)
				return
			}

			log.Println("发送成功!")
		}
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
