package linux

import (
	"fmt"

	"goss.io/goss/lib/command"
	"pandaschool.net/demo/hardware/cmd"
)

//MemSize 内存大小.
func MemSize() string {
	num, err := command.Exec("cat /proc/meminfo | grep 'MemTotal' | awk -F ':' '{print $2}'")
	if err != nil {
		num = "0"
	}

	size := (cmd.ParseInt(num) / (1000 * 1000))
	mem := fmt.Sprintf("%dG", size)

	return mem
}
