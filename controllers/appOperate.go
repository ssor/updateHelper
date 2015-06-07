package controllers

import (
	// "bufio"
	// "fmt"
	// "github.com/Bluek404/downloader"
	// "github.com/astaxie/beego"
	// "github.com/astaxie/beego/config"
	// "github.com/codegangsta/cli"
	"io"
	"os"
	// "strings"
	// "sync"
	"encoding/json"
	// "errors"
	"net/http"
	"os/exec"
	"time"
	// "path"
	// "path/filepath"
)

var (
	chanMonitor      chan bool   = nil
	G_UpdatedAppProc *os.Process = nil
)

//检查是否满足启动条件，满足才能启动
func tryToStartApp() {
	file, err := os.Open(G_appVersionFilePath)
	if err != nil {
		//文件不存在
		DebugMustF("启动应用失败，没有版本信息文件")
		G_currentUpdateInfo = &UpdateInfo{Version: "0.0", FileList: FileChecksumList{}}
		return
	} else {
		bytes := []byte{}
		var buf = make([]byte, 1024)
		for {
			n, e := file.Read(buf)
			bytes = append(bytes, buf[:n]...)
			if e != nil {
				if e == io.EOF {
					// DebugInfo(string(bytes))
					// 数据已经读取完毕
					var cmd UpdateInfo
					if err := json.Unmarshal(bytes, &cmd); err != nil {
						DebugInfoF("读取应用的版本出错:" + err.Error())
					} else {
						// currentVersion = cmd.Version
						G_currentUpdateInfo = &cmd
						// G_currentUpdateInfo.Print()
						DebugInfoF("读取应用的版本: %s", G_currentUpdateInfo.Version)
					}
					break
				}
				break
			}
		}
		file.Close()
		// DebugMust("暂时屏蔽启动应用环节" + GetFileLocation())
		// return
		go TickForStartApp(3)
		go StartAppMonitoring(15)
	}
}

func StartAppMonitoring(second int) {
	if chanMonitor == nil {
		chanMonitor = make(chan bool)
	} else {
		DebugSysF("应用监听已经启动")
		return
	}
	for {
		select {
		case b := <-chanMonitor:
			if b == true {
				DebugInfoF("应用启动成功")
			} else {
				//应用启动失败或者异常退出
				go TickForStartApp(second)
			}
		}
	}
}

//倒数N秒后启动应用
func TickForStartApp(second int) {
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			if second > 1 {
				second = second - 1
				DebugInfoF("应用将在 %d 秒后启动...", second)
			} else {
				StartApp(chanMonitor)
				return
			}
		}
	}
}

func KillApp() error {
	if G_UpdatedAppProc != nil {
		return G_UpdatedAppProc.Kill()
	}
	return nil
}
func StartApp(monitor chan bool) {
	DebugSys("正在启动，大约需要几秒钟...")

	fi, err := os.Stat(G_appBinPath + G_conf.UpdatedAppName)
	if err != nil {
		DebugMustF(err.Error())
		monitor <- false
		return
	}
	if fi.Mode() != os.ModePerm {
		DebugSysF("指定执行文件缺少权限，将尝试提升权限")
		err = os.Chmod(G_appBinPath+G_conf.UpdatedAppName, os.ModePerm)
		if err != nil {
			DebugMustF(err.Error())
			monitor <- false
			return
		}
	}
	var output = &OutputTemp{""}
	go func(binPath string) {
		cmd := exec.Command(binPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = output //错误信息输出
		// cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			DebugMustF("启动应用失败: %s", err.Error())
			monitor <- false
			return
		}
		G_UpdatedAppProc = cmd.Process
		monitor <- true
		if err := cmd.Wait(); err != nil {
			if err.Error() == "signal: killed" {
				DebugInfoF("应用正常退出")
			} else {
				DebugSysF("监控程序不正常退出：%s", err)
				// DebugSysF(output.String())
				// EmailErrorLog("111", output.String(), nil)
				monitor <- false
			}
		}
	}(G_appBinPath + G_conf.UpdatedAppName)
}

func TestAlive() bool {
	resp, err := http.Get("http://localhost:" + G_conf.UpdatedAppPort + "/TestAlive")
	// fmt.Println(resp.Status + "  " + url)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true
}
