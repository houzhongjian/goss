package main

import (
	"log"

	"goss.io/goss/app/master/conf"
	"goss.io/goss/app/master/handler"
	"goss.io/goss/lib/cmd"
)

func main() {
	cmd := cmd.New()

	//加载配置文件.
	conf.Load(cmd)
	log.Println("node name:", conf.Conf.Node.Name)

	master := handler.NewMaster()
	master.Start()
}
