package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func main() {
	st := time.Now().Unix()
	for i := 1; i < 2; i++ {
		upload()
	}

	et := time.Now().Unix()

	log.Println("共耗时:", et-st, "秒")

}

func upload() {
	filename := "./timg.jpeg"
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	writer, err := bodyWriter.CreateFormFile("file", "time.jpg")
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}
	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(writer, fh)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err := http.NewRequest("PUT", "http://127.0.0.1/oss/time.jpg", bodyBuf)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}

	client := http.Client{}
	req.Header.Add("Content-Type", contentType)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}

	log.Println(string(b))
}
