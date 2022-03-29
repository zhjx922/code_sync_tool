package helper

import (
	"log"
	"syscall"
)

func SetLimit() {
	ok := true
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Fatal("获取LIMIT错误：", err)
	}

	rLimit.Cur = rLimit.Max

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Fatal("设置LIMIT错误，手动输入：nlimit -n ", rLimit.Max, "尝试", err)
		ok = false
	}

	if ok {
		log.Println("ulimit 设置成功")
	} else {
		log.Println("ulimit 设置失败")
	}
}
