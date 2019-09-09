package store

import (
	"log"
	"os"

	"github.com/houzhongjian/goini"
)

//Init 初始化存储.
func Init() error {
	//创建存储路径.
	if err := os.MkdirAll(goini.GetString("store_root"), os.ModePerm); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

