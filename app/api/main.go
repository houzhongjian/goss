package main

import (
	"log"

	"goss.io/goss/app/api/conf"
	"goss.io/goss/app/api/handler"
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

	apiSrv := handler.NewApi()
	apiSrv.Start()
}
