package packet

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"goss.io/goss/lib/logd"
	"goss.io/goss/lib/protocol"
)

const PROROCOL_LEN = 4
const HEADER_LEN = 4
const HASH_LEN = 32

type Packet struct {
	Protocol protocol.GossProtocol
	Size     int64
	Hash     string
	Body     []byte
}

type NodeTypes string

const (
	NodeTypes_Master NodeTypes = "master"
	NodeTypes_Store            = "store"
	NodeTypes_Zoo              = "zoo"
)

//NodePacket 节点管理数据.
type NodePacket struct {
	Protocol protocol.GossProtocol //协议号.
	Types    NodeTypes             //节点类型.
	IP       string                //节点ip.
	Size     int64
}

//NewNode.
func NewNode(types NodeTypes, ip string, proto protocol.GossProtocol) []byte {
	body := fmt.Sprintf("%s,%s", types, ip)
	buffer := make([]byte, HEADER_LEN+len(body)+PROROCOL_LEN)
	//0-4 为协议号.
	//4-8 为数据大小.
	//>8 为节点数据.
	binary.BigEndian.PutUint32(buffer[0:4], uint32(proto))
	binary.BigEndian.PutUint32(buffer[4:8], uint32(len(body)))
	copy(buffer[8:], body)
	return buffer
}

func ParseNode(conn net.Conn) (pkt NodePacket, err error) {
	pkt = NodePacket{}
	//获取协议号.
	var num = make([]byte, PROROCOL_LEN)
	_, err = io.ReadFull(conn, num)
	if err != nil {
		log.Printf("%+v\n", err)
		return pkt, err
	}
	pkt.Protocol = protocol.GossProtocol(int(binary.BigEndian.Uint32(num)))

	//获取数据长度.
	var bodyBuf = make([]byte, 4)
	_, err = io.ReadFull(conn, bodyBuf)
	if err != nil {
		log.Printf("%+v\n", err)
		return pkt, err
	}
	pkt.Size = int64(binary.BigEndian.Uint32(bodyBuf))

	//获取数据内容.
	var body = make([]byte, pkt.Size)
	_, err = io.ReadFull(conn, body)
	if err != nil {
		log.Printf("%+v\n", err)
		return pkt, err
	}

	//拆分字符串.
	sArr := strings.Split(string(body), ",")
	pkt.Types = NodeTypes(sArr[0])
	pkt.IP = sArr[1]

	return pkt, nil
}

func New(content, fileHash []byte, num protocol.GossProtocol) []byte {
	buffer := make([]byte, HEADER_LEN+len(content)+len(fileHash)+PROROCOL_LEN)
	//0-4 为协议号.
	//4-8 为文件大小.
	//8-40 为文件hash.
	//>40 为文件内容.
	binary.BigEndian.PutUint32(buffer[0:4], uint32(num))
	binary.BigEndian.PutUint32(buffer[4:8], uint32(len(content)))
	copy(buffer[8:40], fileHash)
	copy(buffer[40:], content)
	return buffer
}

//Parse 解析网络数据包.
func Parse(conn net.Conn) (pkt Packet, err error) {
	pkt = Packet{}
	//获取协议号.
	var num = make([]byte, PROROCOL_LEN)
	_, err = io.ReadFull(conn, num)
	if err != nil {
		log.Printf("%+v\n", err)
		return pkt, err
	}
	pkt.Protocol = protocol.GossProtocol(int(binary.BigEndian.Uint32(num)))

	//获取文件长度.
	var header = make([]byte, HEADER_LEN)
	_, err = io.ReadFull(conn, header)
	if err != nil {
		log.Printf("%+v\n", err)
		return pkt, err
	}
	pkt.Size = int64(binary.BigEndian.Uint32(header))
	logd.Make(logd.Level_INFO, logd.GetLogpath(), fmt.Sprintf("文件大小为:%d", pkt.Size))

	//获取hash.
	var fhash = make([]byte, HASH_LEN)
	_, err = io.ReadFull(conn, fhash)
	if err != nil {
		log.Printf("%+v\n", err)
		return pkt, err
	}
	pkt.Hash = string(fhash)
	logd.Make(logd.Level_INFO, logd.GetLogpath(), fmt.Sprintf("文件hash为:%s", pkt.Hash))

	//获取文件
	var buf = make([]byte, pkt.Size)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		log.Printf("%+v\n", err)
		return pkt, err
	}
	pkt.Body = buf

	return pkt, nil
}
