package db

import (
	"fmt"
	"log"

	"github.com/houzhongjian/goini"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//Db .
var Db *gorm.DB

//DbConfig .
type DbConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     int
	Charset  string
}

//Init .
func Init() error {
	cf := DbConfig{
		Host:     goini.GetString("db_host"),
		User:     goini.GetString("db_user"),
		Password: goini.GetString("db_pwd"),
		Name:     goini.GetString("db_name"),
		Port:     goini.GetInt("db_port"),
		Charset:  goini.GetString("db_charset"),
	}

	return conndb(cf)
}

//conndb .
func conndb(cf DbConfig) error {
	args := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		cf.User,
		cf.Password,
		cf.Host,
		cf.Port,
		cf.Name,
		cf.Charset,
	)
	db, err := gorm.Open("mysql", args)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	db.SingularTable(true)
	db.LogMode(true)

	autoMigrate(db)
	Db = db

	return nil
}

//autoMigrate .
func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&Oss{},
	)
}
