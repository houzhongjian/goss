package main

import (
	"log"

	"goss.io/goss/app/storage/handler"

	"goss.io/goss/app/storage/conf"
	"goss.io/goss/lib/cmd"
)

func main() {
	cmd := cmd.New()

	//加载配置文件.
	conf.Load(cmd)
	log.Println("node name:", conf.Conf.Node.Name)

	storage := handler.NewStorageService()
	storage.Start()
}
