package main

import (
	"log"

	"goss.io/goss/app/master/conf"
	"goss.io/goss/app/master/node"
	"goss.io/goss/db"
	"goss.io/goss/lib/cmd"
)

func main() {
	cmd := cmd.New()

	//加载配置文件.
	conf.Load(cmd)
	log.Println("node name:", conf.Conf.Node.Name)

	if err := db.Connection(); err != nil {
		log.Panicln(err)
	}

	master := node.NewMaster()
	master.Start()
}
