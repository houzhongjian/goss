package db

import (
	"log"

	"github.com/jinzhu/gorm"
)

//Oss .
type Oss struct {
	Model
	Name  string `json:"name"`  //文件名称.
	Store string `json:"store"` //文件存储路径。
	Hash  string `json:"hash"`  //文件hash.
	Size  int64  `json:"size"`  //文件大小.
}

//TableName .
func (o Oss) TableName() string {
	return "oss"
}

//Create.
func (o *Oss) Create() error {
	return Db.Create(o).Error
}

//Query .
func (o *Oss) Query(fname string) (err error) {
	err = Db.Where("name = ?", fname).Order("id desc").First(&o).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("%+v\n", err)
		return err
	}

	return nil
}
