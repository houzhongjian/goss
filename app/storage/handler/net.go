package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

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
