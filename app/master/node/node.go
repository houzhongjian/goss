package node

import (
	"fmt"

	"goss.io/goss/app/master/api"
	"goss.io/goss/app/master/conf"
)

//Gateway 网关.
type Gateway struct {
}

func NewMaster() *Gateway {
	gwy := &Gateway{}
	return gwy
}

//Start .
func (g *Gateway) Start() {
	cf := conf.Conf.Node
	restartMsg := fmt.Sprintf("%s(%s:%d) 启动成功!", cf.Name, cf.IP, cf.Port)
	println(restartMsg)
	//创建对外api接口.
	g.newApi()
	select {}
}

//newApi .
func (g *Gateway) newApi() {
	api.NewService()
}
