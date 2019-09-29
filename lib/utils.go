package lib

import (
	"crypto/md5"
	"encoding/hex"
)

func FileHash(body []byte) string {
	h := md5.New()
	h.Write(body)
	b := h.Sum(nil)
	return hex.EncodeToString(b)
}
