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
	G_versionUrl                                     = "https://raw.githubusercontent.com/ssor/binpickup/master/VersionInfo.md"
	G_baseUrl                                        = "https://raw.githubusercontent.com/ssor/binpickup/master/"
	G_UpdatedAppProc          *os.Process            = nil
	G_errNotAllTaskCompleted  error                  = errors.New("有未完成的下载任务")
	G_downloadingUpdateInfo   *UpdateInfo            = nil //当前运行的应用的版本信息
	G_downloadTasks           DownloadTaskList       = DownloadTaskList{}
)

var ( //应用的相关信息
	G_currentUpdateInfo  *UpdateInfo = nil                    //当前运行的应用的版本信息
	G_UpdatedAppPort     string      = ""                     //http端口
	G_UpdatedAppName     string      = ""                     //应用的名称
	G_appBasePath                    = "./App/"               //应用所在的目录
	G_appBinPath                     = "./App/Bin/"           //应用所在的目录
	G_appVersionFilePath             = "./App/VersionInfo.md" //应用的版本文件所在位置
	G_UpdateResourcePath             = "./UpdateResource/"
	G_versionInfoFile                = "VersionInfo.md"
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

//----------------------------------------------------------------------
