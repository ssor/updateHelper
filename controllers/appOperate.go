package controllers

import (
	// "bufio"
	"fmt"
	// "github.com/Bluek404/downloader"
	// "github.com/astaxie/beego"
	// "github.com/astaxie/beego/config"
	// "github.com/codegangsta/cli"
	"io"
	"os"
	// "strings"
	// "sync"
	"encoding/json"
	"errors"
	"net/http"
	"os/exec"
	"time"
	// "path"
	// "path/filepath"
)

//检查是否满足启动条件，满足才能启动
func tryToStartApp() {
	file, err := os.Open(G_appVersionFilePath)
	if err != nil {
		//文件不存在
		DebugMust("启动应用失败，没有版本信息文件" + GetFileLocation())
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
						DebugInfo("读取应用的版本出错:" + err.Error() + GetFileLocation())
					} else {
						// currentVersion = cmd.Version
						G_currentUpdateInfo = &cmd
						// G_currentUpdateInfo.Print()
						DebugInfo("读取应用的版本: " + G_currentUpdateInfo.Version + GetFileLocation())
					}
					break
				}
				break
			}
		}
		file.Close()
		DebugMust("暂时屏蔽启动应用环节")
		return
		if StartApp() == true {
			DebugInfo("启动应用成功" + GetFileLocation())
		} else {
			DebugInfo("启动应用失败" + GetFileLocation())
		}
	}
}
func KillApp() error {
	if G_UpdatedAppProc != nil {
		return G_UpdatedAppProc.Kill()
	}
	if TestAlive() == true {
		return errors.New("应用非正常启动，需要手动关闭应用")
	} else {
		return errors.New("应用没有启动")
	}
}
func StartApp() bool {
	go func(binPath string) {
		cmd := exec.Command(binPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// log.Print(cmdline)
		err := cmd.Start()
		if err != nil {
			DebugMust(fmt.Sprintf("启动应用失败: '%s'\n", err))
		}
		G_UpdatedAppProc = cmd.Process
	}(G_appBasePath + G_UpdatedAppName)
	DebugSys("正在启动，大约需要几秒钟。。。")
	time.Sleep(time.Second * 3)
	return TestAlive()
	// return true
}

func TestAlive() bool {
	resp, err := http.Get("http://localhost:" + G_UpdatedAppPort + "/TestAlive")
	// fmt.Println(resp.Status + "  " + url)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	// bytes := []byte{}
	// var buf = make([]byte, 1024)
	// for {
	// 	n, e := resp.Body.Read(buf)
	// 	bytes = append(bytes, buf[:n]...)
	// 	if e != nil {
	// 		if e == io.EOF {
	// 			// fmt.Println("数据读取完毕")
	// 			// 数据已经下载完毕
	// 			var cmd Command
	// 			if err := json.Unmarshal(bytes, &cmd); err != nil {
	// 				DebugInfo("获取应用的版本出错")
	// 			} else {
	// 				currentVersion = cmd.Message
	// 				DebugInfo("获取应用的版本: " + currentVersion)
	// 			}
	// 			break
	// 		}
	// 		break
	// 	}
	// }
	return true
}
