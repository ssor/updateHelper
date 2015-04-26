package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"runtime"
	"strings"
)

/*
	开发系统的核心思想是通过一个模型对数据依据一定的逻辑规则进行处理，模型应该在相应环节发出系统运行提示，可以分为4级
	1  必须处理，否则系统无法正常运转
	2  可以稍后处理的异常，影响系统运行的提示信息，例如参数变化，外力因素导致
	3  提示信息，内部数据变化
	4  调试诊断用，生产环境不会使用
	相应可以添加更低级别的提示，以适应调试的需要
*/

var (
	LevelTrace    = 6 //运行数据
	LevelDebug    = 5 //运行数据
	LevelInfo     = 4 //运行数据
	LevelWarn     = 3 //处理错误
	LevelError    = 2 //处理错误
	LevelCritical = 1 //处理错误
)

func init() {
	beego.BeeLogger.EnableFuncCallDepth(false)
	beego.BeeLogger.SetLogFuncCallDepth(7)

	// beego.Trace("traceeeeeee")
	// beego.Debug("debuggggggg")
	// beego.Info("infooooooooo")
	// beego.Warn("warnnnnnnnnn")
	// beego.Error("errorrrrrr")
	// beego.Critical("criticalllllll")
	//"test"

}

var DebugLevel int = 4

var userBeego = false

var G_printLog = true

// var userBeego = true

func DebugMust(log string) {
	DebugOutput(log, 1)
}
func DebugSys(log string) {
	DebugOutput(log, 2)
}
func DebugInfo(log string) {
	DebugOutput(log, 3)
}
func DebugTrace(log string) {
	DebugOutput(log, 4)
}

func DebugOutput(log string, level int) {
	if G_printLog == false {
		return
	}
	if level <= DebugLevel {
		if userBeego == true {
			DebugOutputBeego(log, level)
		} else {

			prefix := ""
			switch level {
			case 1:
				prefix = "#  ******"
			case 2:
				prefix = "#  ------"
			case 3:
				prefix = "#        "
			case 4:
				prefix = "#               "

			}
			fmt.Println(prefix + log)
		}
	}
}

func DebugOutputBeego(log string, level int) {
	switch level {
	case 1:
		beego.Error(log)
	case 2:
		beego.Notice(log)
	case 3:
		beego.Informational(log)
	case 4:
		beego.Debug(log)
	}
}
func GetFileLocation() string {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		array := strings.Split(file, "/")
		return fmt.Sprintf(" (%s %d)", array[len(array)-1], line)
	} else {
		return "  ???"
	}
}

func DebugOutputStrings(strs []string, level int) {
	for _, str := range strs {
		DebugOutput(str, level)
	}
}
