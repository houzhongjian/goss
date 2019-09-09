package service

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/houzhongjian/goini"
	"pandaschool.net/goss/src/db"
	"pandaschool.net/goss/src/utils"
)

//Handler .
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		put(w, r)
		return
	}

	if r.Method == http.MethodGet {
		get(w, r)
		return
	}

	if r.Method == http.MethodDelete {
		del(w, r)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

//put.
func put(w http.ResponseWriter, r *http.Request) {
	fname := strings.Split(r.URL.EscapedPath(), "/")[2]
	// ftype := strings.Split(r.URL.EscapedPath(), "/")[3]
	log.Println(fname)
	fpath := utils.FileStorePath()
	if err := os.MkdirAll(goini.GetString("store_root")+fpath, os.ModePerm); err != nil {
		log.Printf("创建文件夹失败:%+v\n", err)
		w.Write([]byte("FAIL"))
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("copy文件内容失败:%+v\n", err)
		w.Write([]byte("FAIL"))
		return
	}

	storePath, hash, err := fileExist(b)
	if err != nil {
		log.Printf("copy文件内容失败:%+v\n", err)
		w.Write([]byte("FAIL"))
		return
	}
	oss := db.Oss{
		Name:  fname,
		Store: fpath,
		Hash:  hash,
		// Type:  ftype,
		Size: r.ContentLength,
	}

	if len(storePath) < 1 {
		f, err := os.Create(goini.GetString("store_root") + fpath + fname)
		if err != nil {
			log.Printf("创建文件失败:%+v\n", err)
			w.Write([]byte("FAIL"))
			return
		}
		defer f.Close()

		_, err = f.Write(b)
		if err != nil {
			log.Printf("copy文件内容失败:%+v\n", err)
			w.Write([]byte("FAIL"))
			return
		}
		oss.Store = fpath + fname
	}

	if err = oss.Create(); err != nil {
		log.Printf("文件记录写入数据库失败%+v\n", err)
		w.Write([]byte("FAIL"))
		return
	}

	w.Write([]byte("SUCCESS"))
}

//get.
func get(w http.ResponseWriter, r *http.Request) {
	fname := strings.Split(r.URL.EscapedPath(), "/")[2]
	//根据文件名获取文件存储的路径.
	oss := db.Oss{}
	if err := oss.Query(fname); err != nil {
		log.Printf("获取文件路径失败:%+v\n", err)
		w.Write([]byte("FAIL"))
		return
	}

	if len(oss.Store) < 1 {
		w.Write([]byte("not found"))
		return
	}

	store := goini.GetString("store_root")
	fpath := store + oss.Store
	f, err := os.Open(fpath)
	if err != nil {
		log.Printf("%+v\n", err)
		w.Write([]byte("FAIL"))
		return
	}

	io.Copy(w, f)
}

//del .
func del(w http.ResponseWriter, r *http.Request) {

}

//fileExist 验证文件是否存在.
func fileExist(b []byte) (storePath, hash string, err error) {
	hash = utils.FileMD5(b)

	oss := db.Oss{}
	err = db.Db.Model(&oss).Where("hash = ?", hash).First(&oss).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("%+v\n", err)
		return "", hash, err
	}

	return oss.Store, hash, nil
}
