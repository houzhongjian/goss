package main

import (
	"log"

	"goss.io/goss/app/store/handler"

	"goss.io/goss/app/store/conf"
	"goss.io/goss/lib/cmd"
)

func main() {
	cmd := cmd.New()

	//加载配置文件.
	conf.Load(cmd)
	log.Println("node name:", conf.Conf.Node.Name)

	store := handler.NewStoreService()
	store.Start()
}
