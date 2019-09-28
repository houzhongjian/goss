package cmd

import (
	"flag"
	"log"
	"os"
)

func (this *Command) parse() {

	var confPath string
	flag.StringVar(&confPath, "conf", "", "conf")

	flag.Parse()

	//conf.
	this.Conf = confPath
	if confPath == "" {
		log.Println("请指定配置文件 -conf")
		os.Exit(0)
	}
}

type Command struct {
	Conf string
}

func New() *Command {
	command := &Command{}
	command.parse()
	command.helloView()

	return command
}

func (this *Command) helloView() {
	var cmdMsg = `Goss Version 1.0`
	println(cmdMsg)
}
