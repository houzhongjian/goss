package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"goss.io/goss/db"
	"goss.io/goss/lib/filetype"
	"goss.io/goss/lib/protocol"

	"goss.io/goss/app/master/conf"
	"goss.io/goss/app/master/handler"
	"goss.io/goss/lib"
	"goss.io/goss/lib/packet"
)

//ApiService.
type ApiService struct {
	Port string
	Tcp  *handler.TcpService
}

//NewService .
func NewService() {
	cf := conf.Conf.Node
	apiSrv := ApiService{
		Port: fmt.Sprintf(":%d", cf.Port),
		Tcp:  handler.NewTcpService(),
	}
	go apiSrv.Tcp.Start()
	apiSrv.httpSrv()
}

//httpSrv .
func (this *ApiService) httpSrv() {
	http.HandleFunc("/oss/", this.handler)
	if err := http.ListenAndServe(this.Port, nil); err != nil {
		log.Panicf("%+v\n", err)
	}
}

func (this *ApiService) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		this.get(w, r)
		return
	}

	if r.Method == http.MethodPut {
		this.put(w, r)
		return
	}

	if r.Method == http.MethodDelete {
		this.delete(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

//get.
func (this *ApiService) get(w http.ResponseWriter, r *http.Request) {
	name, err := this.getParse(r.URL.EscapedPath())
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	meta := db.Metadata{
		Name: name,
	}
	if err = meta.Query(); err != nil {
		log.Printf("%+v\n", err)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Printf("metadata:%+v\n", meta)

	b, err := this.Tcp.Read(meta.StoreNode, meta.Hash, meta.Size)
	if err != nil {
		log.Printf("%+v\n", err)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte(b))
}

//getParse get请求解析文件名.
func (this *ApiService) getParse(url string) (name string, err error) {
	sArr := strings.Split(url, "/")
	if len(sArr) != 3 {
		return name, errors.New("not fount")
	}
	if sArr[2] == "" {
		return name, errors.New("not fount")
	}

	return sArr[2], nil
}

//put.
func (this *ApiService) put(w http.ResponseWriter, r *http.Request) {
	//获取文件名称，文件大小，文件类型，文件hash.
	//元数据.
	name, err := this.getParse(r.URL.EscapedPath())
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%+v\n", err)
		w.Write([]byte("fail"))
		return
	}

	//获取文件类型.
	f16 := fmt.Sprintf("%x", fBody)
	ft := filetype.Parse(f16[:10])

	fhash := lib.FileHash(fBody)
	pkt := packet.New(fBody, []byte(fhash), protocol.WriteFileProrocol)
	_, nodeip, err := this.Tcp.Write(pkt)
	if err != nil {
		log.Printf("%+v\n", err)
		w.Write([]byte("fail"))
		return
	}

	//记录文件元数据.
	metadata := db.Metadata{
		Name:      name,
		Type:      ft,
		Size:      int64(len(fBody)),
		Hash:      fhash,
		StoreNode: nodeip,
	}

	if err := metadata.Create(); err != nil {
		log.Printf("%+v\n", err)
		w.Write([]byte("fail"))
		return
	}

	w.Write([]byte("success"))
}

//delete.
func (this *ApiService) delete(w http.ResponseWriter, r *http.Request) {

}
