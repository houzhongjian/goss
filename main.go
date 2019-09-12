package main

import (
	"log"

	"pandaschool.net/goss/src/cmd"
	"pandaschool.net/goss/src/conf"
	"pandaschool.net/goss/src/node"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	cmd := cmd.New()

	//加载配置文件.
	conf.Load(cmd)
	log.Println("node name:", conf.Conf.Node.Name)

	//判断是否为主节点.
	if cmd.MasterNode {
		log.Println("当前节点为:主节点")
		//获取存储节点.
		g := node.NewMaster()
		g.Start()
		return
	}

	//创建存储节点.
	s := node.NewStore()
	s.Start()
}
