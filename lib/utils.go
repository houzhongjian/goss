package lib

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"time"
)

func FileHash(body []byte) string {
	h := md5.New()
	h.Write(body)
	b := h.Sum(nil)
	return hex.EncodeToString(b)
}

//IsExists 判断ini是否存在.
func IsExists(ini string) bool {
	_, err := os.Stat(ini)
	if err != nil {
		return false
	}
	return true
}

//Time 获取当前时间.
func Time() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
