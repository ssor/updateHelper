package controllers

import (
	// "bufio"
	// "fmt"
	// "github.com/Bluek404/downloader"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	// "github.com/codegangsta/cli"
	// "io"
	"os"
	// "strings"
	// "sync"
	// "encoding/json"
	"errors"
	// "net/http"
	// "os/exec"
	// "time"
	// "path"
	// "path/filepath"
)

var (
	G_CheckUpdateIntervalMode bool                   = true
	G_iniconf                 config.ConfigContainer = nil
	G_versionInfoFile                                = "VersionInfo.md"
	G_baseUrl                                        = "https://raw.githubusercontent.com/ssor/binpickup/master/"
	G_versionUrl                                     = "https://raw.githubusercontent.com/ssor/binpickup/master/" + G_versionInfoFile
	G_UpdatedAppProc          *os.Process            = nil
	G_errNotAllTaskCompleted  error                  = errors.New("有未完成的下载任务")
	G_downloadingUpdateInfo   *UpdateInfo            = nil //当前运行的应用的版本信息
	G_downloadTasks           DownloadTaskList       = DownloadTaskList{}
	G_updateAppReady          bool                   = false
)

var ( //应用的相关信息
	G_currentUpdateInfo  *UpdateInfo = nil                          //当前运行的应用的版本信息
	G_UpdatedAppPort     string      = ""                           //http端口
	G_UpdatedAppName     string      = ""                           //应用的名称，实际的执行文件
	G_appBasePath                    = "./App/"                     //应用所在的目录
	G_appBinPath                     = "./App/Bin/"                 //应用所在的目录
	G_appVersionFilePath             = "./App/" + G_versionInfoFile //应用的版本文件所在位置
	G_UpdateResourcePath             = "./UpdateResource/"
	// G_appID              string      = ""                           //与应用通信的标识，一般是应用的名称
)

type Command struct {
	Code    int
	Message string
	Data    interface{}
}

func NewCommandData(code int, msg string, data interface{}) Command {
	cmd := newCommand(code, msg)
	cmd.Data = data
	return cmd
}
func newCommand(code int, msg string) Command {
	return Command{
		Code:    code,
		Message: msg,
	}
}

type MainController struct {
	beego.Controller
}

func (this *MainController) TestAlive() {
	this.Data["json"] = newCommand(0, beego.AppName)
	this.ServeJson()
}

//呼叫升级助手可以进行升级了，应用根据返回值确定升级助手现在是否方便升级
func (this *MainController) StartUpdate() {
	// appName := this.GetString("App")
	// this.Data["json"] = newCommand(1, "应用名称不对应")
	// DebugInfoF("应用名称不对应，传入的应用名称为：%s", appName)
	// if appName == G_appID {
	go UpdateApp()
	this.Data["json"] = newCommand(0, "")
	// }
	this.ServeJson()
}

//----------------------------------------------------------------------
