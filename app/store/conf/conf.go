package conf

import (
	"log"
	"os"

	"goss.io/goss/lib"

	"goss.io/goss/lib/cmd"
	"goss.io/goss/lib/ini"
)

type Config struct {
	Node *nodeConfig
	Base *baseConfig
}

type nodeConfig struct {
	IP        string
	Port      int
	Name      string
	StoreRoot string
	ZooNode   string
}

type baseConfig struct {
	LogPath string
}

var Conf *Config

//Load .
func Load(cmd *cmd.Command) {
	iniPath := cmd.Conf
	if !lib.IsExists(cmd.Conf) {
		log.Println("配置文件不存在=>", iniPath)
		os.Exit(0)
		return
	}

	if err := ini.Load(iniPath); err != nil {
		log.Printf("%+v\n", err)
		return
	}

	cf := &Config{
		Node: parseNodeConfig(cmd),
		Base: parseBaseConfig(),
	}

	Conf = cf
}

//node.
func parseNodeConfig(cmd *cmd.Command) *nodeConfig {
	name := ini.GetString("node_name")
	if len(name) < 1 {
		log.Println("node_name 不能为空")
		os.Exit(0)
	}

	storeip := ini.GetString("node_ip")
	if len(storeip) < 1 {
		log.Println("node_ip 不能为空")
		os.Exit(0)
	}
	storeport := ini.GetInt("node_port")
	if storeport < 1 {
		log.Println("node_port 不能为空")
		os.Exit(0)
	}

	storeRoot := ini.GetString("store_root")
	if len(storeRoot) < 1 {
		log.Println("store_root 不能为空")
		os.Exit(0)
	}
	zooNode := ini.GetString("zoo_node")
	if len(storeRoot) < 1 {
		log.Println("zoo_node 不能为空")
		os.Exit(0)
	}

	nodeconf := &nodeConfig{
		IP:        storeip,
		Port:      storeport,
		Name:      name,
		StoreRoot: storeRoot,
		ZooNode:   zooNode,
	}

	return nodeconf
}

//parseBaseConfig 解析基础配置.
func parseBaseConfig() *baseConfig {
	logpath := ini.GetString("log_path")
	if len(logpath) < 1 {
		log.Println("log_path 不能为空")
		os.Exit(0)
	}
	base := baseConfig{
		LogPath: logpath,
	}
	return &base
}
