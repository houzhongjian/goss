package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

//FileStorePath 生成文件存储路径.
func FileStorePath() string {
	t := time.Now()

	year := fmt.Sprintf("%d", t.Year())
	month := fmt.Sprintf("%d", t.Month())
	day := fmt.Sprintf("%d", t.Day())
	hour := fmt.Sprintf("%d", t.Hour())
	minute := fmt.Sprintf("%d", t.Minute())

	return fmt.Sprintf("%s/%s/%s/%s/%s/", year, month, day, hour, minute)
}

//FileMD5 判读文件是否已经存在.
func FileMD5(b []byte) string {
	h := md5.New()
	h.Write(b)
	b = h.Sum(nil)
	return hex.EncodeToString(b)
}
