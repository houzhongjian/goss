package conf

import (
	"log"
	"os"

	"github.com/houzhongjian/goini"
	"goss.io/goss/app/store/src/cmd"
)

type Config struct {
	Node *nodeConfig
	Db   *dbConfig
}

type dbConfig struct {
	Host     string
	User     string
	Port     int
	Name     string
	Password string
}

type nodeConfig struct {
	IP   string
	Port int
	Name string
}

var Conf *Config

//Load .
func Load(cmd *cmd.Command) {
	ini := cmd.Conf
	if !iniIsExists(cmd.Conf) {
		log.Println("配置文件不存在=>", ini)
		os.Exit(0)
		return
	}

	if err := goini.Load(ini); err != nil {
		log.Printf("%+v\n", err)
		return
	}

	cf := &Config{
		Db:   parseDbConfig(),
		Node: parseNodeConfig(cmd),
	}

	Conf = cf
}

//node.
func parseNodeConfig(cmd *cmd.Command) *nodeConfig {
	name := goini.GetString("node_name")
	if len(name) < 1 {
		log.Println("node_name 不能为空")
		os.Exit(0)
	}

	storeip := goini.GetString("node_ip")
	if len(storeip) < 1 {
		log.Println("node_ip 不能为空")
		os.Exit(0)
	}
	storeport := goini.GetInt("node_port")
	if storeport < 1 {
		log.Println("node_port 不能为空")
		os.Exit(0)
	}

	nodeconf := &nodeConfig{
		IP:   storeip,
		Port: storeport,
		Name: name,
	}

	return nodeconf
}

//db.
func parseDbConfig() *dbConfig {
	dbHost := goini.GetString("db_host")
	if len(dbHost) < 1 {
		log.Println("db_host 不能为空")
		os.Exit(0)
	}

	dbUser := goini.GetString("db_user")
	if len(dbUser) < 1 {
		log.Println("db_user 不能为空")
		os.Exit(0)
	}

	dbPort := goini.GetInt("db_port")
	if dbPort < 1 {
		log.Println("db_port 不能为空")
		os.Exit(0)
	}

	dbName := goini.GetString("db_name")
	if len(dbName) < 1 {
		log.Println("db_name 不能为空")
		os.Exit(0)
	}

	dbPwd := goini.GetString("db_pwd")
	if len(dbPwd) < 1 {
		log.Println("db_pwd 不能为空")
		os.Exit(0)
	}

	dbconf := &dbConfig{
		Host:     dbHost,
		User:     dbUser,
		Port:     dbPort,
		Name:     dbName,
		Password: dbPwd,
	}
	return dbconf
}

//iniIsExists 判断ini是否存在.
func iniIsExists(ini string) bool {
	_, err := os.Stat(ini)
	if err != nil {
		return false
	}
	return true
}
