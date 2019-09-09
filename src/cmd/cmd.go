package cmd

import (
	"flag"
	"log"
	"os"
)

//RunMode .
type RUN_MODE string

const (
	RUN_MODE_DEV  RUN_MODE = "dev"
	RUN_MODE_PROD          = "prod"
)

//NODE_TYPE.
type NODE_TYPE string

const (
	NODE_TYPE_MASTER NODE_TYPE = "master"
	NODE_TYPE_STORE            = "store"
)

func (c *Command) parse() {
	var nodeType string
	flag.StringVar(&nodeType, "node", "store", "node")

	var confPath string
	flag.StringVar(&confPath, "conf", "", "conf")

	flag.Parse()

	//nodetype.
	c.NodeType = NODE_TYPE_STORE
	if nodeType == "master" {
		c.NodeType = NODE_TYPE_MASTER
		c.MasterNode = true
	}

	//conf.
	c.Conf = confPath
	if confPath == "" {
		log.Println("请指定配置文件 -conf")
		os.Exit(0)
	}
}

type Command struct {
	NodeType   NODE_TYPE
	MasterNode bool
	Conf       string
}

func New() *Command {
	command := &Command{}
	command.parse()

	return command
}
