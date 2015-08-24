package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	// "github.com/fatih/color"
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

	if strings.Contains(strings.ToLower(runtime.GOOS), "windows") == true { //windows平台下不使用beego打印log的方案
		useBeego = false
	}
}

type IPrintList interface {
	ListName() string
	InfoList() []string
}

var DebugLevel int = 4

var useBeego = true

// var useColor = false

var G_printLog = true
var G_DebugLine = "-------------------------------------------------------------------------"

func DebugPrintList_Info(list IPrintList) {
	PrintList(list, DebugInfo)
}
func DebugPrintList_Trace(list IPrintList) {
	PrintList(list, DebugTrace)
}
func PrintList(list IPrintList, printFunc func(log string)) {
	log := fmt.Sprintf(G_DebugLine+"%s", getFileLocation())
	printFunc(log)
	printFunc(list.ListName() + " 列表：")
	strs := list.InfoList()
	for _, str := range strs {
		printFunc(str)
	}
	printFunc(log)
}

//能够造成系统不正常运行的问题
func DebugMust(log string) {
	DebugOutput(log, 1)
}
func DebugMustF(format string, args ...interface{}) {
	log := fmt.Sprintf(format+"%s", append(args, getFileLocation())...)
	DebugMust(log)
}
func DebugSysF(format string, args ...interface{}) {
	log := fmt.Sprintf(format+"%s", append(args, getFileLocation())...)
	DebugSys(log)
}

// 出现异常信息，系统能够正常运行，但是可能和使用者想象的不同
func DebugSys(log string) {
	DebugOutput(log, 2)
}

func DebugInfoF(format string, args ...interface{}) {
	log := fmt.Sprintf(format+"%s", append(args, getFileLocation())...)
	DebugInfo(log)
}

// 关键步骤或者信息的提醒
func DebugInfo(log string) {
	DebugOutput(log, 3)
}

func DebugTraceF(format string, args ...interface{}) {
	log := fmt.Sprintf(format+"%s", append(args, getFileLocation())...)
	DebugTrace(log)
}

// 运行数据的打印
func DebugTrace(log string) {
	DebugOutput(log, 4)
}

func DebugOutput(log string, level int) {
	if G_printLog == false {
		return
	}
	if level <= DebugLevel {
		if useBeego == true {
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
func getFileLocation() string {
	_, file, line, ok := runtime.Caller(2)
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
