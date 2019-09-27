package main

import (
	"log"

	"goss.io/goss/app/master/src/cmd"
	"goss.io/goss/app/master/src/conf"
	"goss.io/goss/app/master/src/node"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	cmd := cmd.New()

	//加载配置文件.
	conf.Load(cmd)
	log.Println("node name:", conf.Conf.Node.Name)

	master := node.NewMaster()
	master.Start()
}
