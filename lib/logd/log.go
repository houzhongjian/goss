package logd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
	"goss.io/goss/lib"
	"goss.io/goss/lib/ini"
)

//Level .
type LogLevel string

const (
	Level_INFO    LogLevel = "INFO"
	Level_DEBUG            = "DEBUG"
	Level_WARNING          = "WARNING"
	Level_ERROR            = "ERROR"
)

func Make(level LogLevel, logpath string, msg interface{}) {
	makelog(level, logpath, msg)
}

func makelog(level LogLevel, logpath string, msg interface{}) {
	nodename := ini.GetString("node_name")
	content := fmt.Sprintf("%s %s:[%s] %s [%v]\n", lib.Time(), nodename, level, logpath, msg)

	logFile := logFile()
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDONLY|os.O_WRONLY, 0777)
	if err != nil {
		log.Printf("%+v\n", err)
		f.Close()
		return
	}
	println(content)
	_, err = f.WriteString(content)
	if err != nil {
		log.Printf("%+v\n", err)
	}
	f.Close()
}

//getLogpath 获取产生日志的路径.
func GetLogpath() string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", file, line)
}

func logFile() string {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	day := now.Day()
	hour := now.Hour()

	path := fmt.Sprintf("%s%d/%d/%d/", ini.GetString("log_path"), year, month, day)

	//判断存储路径是否存在.
	if !lib.IsExists(path) {
		os.MkdirAll(path, 0777)
	}

	return fmt.Sprintf("%s%d.log", path, hour)
}
