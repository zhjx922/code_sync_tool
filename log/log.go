package log

import (
	"github.com/fatih/color"
	"log"
)

func Println(v ...interface{}) {
	//info := color.New(color.FgWhite, color.BgGreen).SprintFunc()
	log.Println(v...)
}

func Printf(format string, a ...interface{}) {
	//info := color.New(color.FgWhite, color.BgGreen).SprintFunc()
	log.Printf(format, a...)
}

func Warning(v ...interface{}) {
	info := color.New(color.FgWhite, color.BgRed).SprintFunc()
	log.Println(info(v...))
}
