package main

import (
	"log"

	"goss.io/goss/app/store/handler"

	"goss.io/goss/app/store/cmd"
	"goss.io/goss/app/store/conf"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	cmd := cmd.New()

	//加载配置文件.
	conf.Load(cmd)
	log.Println("node name:", conf.Conf.Node.Name)

	store := handler.NewStoreService()
	store.Start()
}
