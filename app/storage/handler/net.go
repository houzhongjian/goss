package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

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
		log.Printf("端口监听失败!%+v\n", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil && err == io.EOF {
			logd.Make(logd.Level_WARNING, logd.GetLogpath(), "断开连接")
			return
		}

		ip := conn.RemoteAddr().String()
		go this.handler(conn, ip)
	}
}

func (this *StorageService) handler(conn net.Conn, ip string) {
	defer conn.Close()
	for {
		pkt, err := packet.Parse(conn)
		if err != nil && err == io.EOF {
			logd.Make(logd.Level_WARNING, logd.GetLogpath(), ip+"断开连接")
			return
		}

		//判断协议号.
		if pkt.Protocol == protocol.SEND_FILE {
			//计算文件hash.
			fHash := lib.FileHash(pkt.Body)
			//验证文件是否损坏.
			if fHash != pkt.Hash {
				logd.Make(logd.Level_WARNING, logd.GetLogpath(), "文件hash不一致")
				conn.Write([]byte("fail"))
				return
			}

			fPath := conf.Conf.Node.StorageRoot + fHash
			err = ioutil.WriteFile(fPath, pkt.Body, 0777)
			if err != nil {
				logd.Make(logd.Level_WARNING, logd.GetLogpath(), "创建文件失败"+err.Error())
				return
			}
			conn.Write([]byte(fHash))
		}

		if pkt.Protocol == protocol.READ_FILE {
			//读取文件.
			fpath := conf.Conf.Node.StorageRoot + pkt.Hash
			b, err := ioutil.ReadFile(fpath)
			if err != nil {
				logd.Make(logd.Level_WARNING, logd.GetLogpath(), "读取文件失败:"+err.Error())
				return
			}
			_, err = conn.Write(b)
			if err != nil && err == io.EOF {
				logd.Make(logd.Level_WARNING, logd.GetLogpath(), "文件发送失败:"+err.Error())
				return
			}

			logd.Make(logd.Level_INFO, logd.GetLogpath(), "文件发成功")
		}
	}
}
