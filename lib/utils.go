package lib

import (
	"crypto/md5"
	"encoding/hex"
	"os"
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
