package controllers

import (
	"bufio"
	"fmt"
	// "github.com/Bluek404/downloader"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/codegangsta/cli"
	// "io"
	"os"
	"strings"
	// "sync"
	// "encoding/json"
	// "errors"
	"net/http"
	// "os/exec"
	"time"
	// "path"
	// "path/filepath"
)

func init() {
	// copyUpdateFileToApp()
	// return
	initConfig()
	go initCli()

	tryToStartApp()
	time.Sleep(time.Second * 2)

	go startIntervalCheckUpdateInfoFromServer()
	go startIntervalNotifyAppUpdate()
}
func startIntervalNotifyAppUpdate() {
	c := time.Tick(5 * time.Second)
	// c := time.Tick(1 * time.Minute)
	for range c {
		if G_updateAppReady == true {
			go func() {
				DebugInfo("提示App可以升级了" + GetFileLocation())
				resp, err := http.Get("http://localhost:" + G_UpdatedAppPort + "/Update")
				if err != nil {
					DebugSys(fmt.Sprintf("App无法接收升级提示: %s", err.Error()) + GetFileLocation())
				}
				if resp.StatusCode != 200 {
					DebugSys(fmt.Sprintf("App无法接收升级提示: %s", resp.Status) + GetFileLocation())
				}
			}()
		} else {
			DebugInfo("没有升级信息可以通知App" + GetFileLocation())
		}
	}
}
func startIntervalCheckUpdateInfoFromServer() {
	c := time.Tick(15 * time.Second)
	// c := time.Tick(1 * time.Minute)
	for range c {
		if G_CheckUpdateIntervalMode == true {
			go func() {
				DebugInfo("定时检查升级模式启用" + GetFileLocation())
				CheckUpdate()
			}()
		} else {
			DebugInfo("定时检查升级模式关闭" + GetFileLocation())
		}
	}
}

func initConfig() {
	var err error
	G_iniconf, err = config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		beego.Error(err.Error())
	} else {
		updatedAppPort := G_iniconf.String("updatedAppPort")
		if len(updatedAppPort) <= 0 {
			G_UpdatedAppPort = "9001"
		} else {
			G_UpdatedAppPort = updatedAppPort
		}
		DebugInfo("升级目标应用的端口：" + G_UpdatedAppPort + GetFileLocation())

		G_UpdatedAppName = G_iniconf.String("updatedAppName")
		DebugInfo("升级目标应用的名称：" + G_UpdatedAppName + GetFileLocation())

		// if err := iniconf.Set("locationCount", "23"); err != nil {
		// 	beego.Warn(err.Error())
		// }
		// if err := iniconf.SaveConfigFile("conf/app.conf"); err != nil {
		// 	beego.Warn(err.Error())
		// }
	}
}

func initCli() {
	cliApp := cli.NewApp()
	cliApp.Name = ""
	cliApp.Usage = "设置系统运行参数"
	cliApp.Version = "1.0.1"
	cliApp.Email = "ssor@qq.com"
	cliApp.Commands = []cli.Command{
		{
			Name:        "update",
			ShortName:   "ud",
			Usage:       "系统升级",
			Description: "如果有系统更新，下载升级文件",
			Action: func(c *cli.Context) {
				// fmt.Println(fmt.Sprintf("%#v", c.Command))
				// fmt.Println("-----------------------------")
				value := strings.ToLower(c.Args().First())
				beego.Info(fmt.Sprintf("查看系统更新：%s", value))
				// DownloadFromUrl(urlDir)
				// tasks := DownloadTaskList{NewDownloadTask(url1), NewDownloadTask(url2)}
				// tasks := DownloadTaskList{NewDownloadTask(url1), NewDownloadTask(url2), NewDownloadTask(url3)}
				// StartDownloadTask(tasks)
				CheckUpdate()
			},
		}, {
			Name:        "getUpdateInfo",
			ShortName:   "gu",
			Usage:       "查看是否有更新信息",
			Description: "首先获取版本和升级信息",
			Action: func(c *cli.Context) {
				// fmt.Println(fmt.Sprintf("%#v", c.Command))
				// fmt.Println("-----------------------------")
				// root := strings.ToLower(c.Args().First())
				// if len(root) <= 0 {
				// 	root = "Bin"
				// }
				beego.Info(fmt.Sprintf("获取要升级的版本信息"))
				GetVersionInfo(G_versionUrl)
			},
		}, {
			Name:        "StartApp",
			ShortName:   "sa",
			Usage:       "启动要升级的应用",
			Description: "要启动的应用位于 " + G_appBasePath + " 目录下",
			Action: func(c *cli.Context) {
				if StartApp() == true {
					DebugInfo("启动应用成功")
				} else {
					DebugInfo("启动应用失败")
				}
			},
		}, {
			Name:        "TestAlive",
			ShortName:   "ta",
			Usage:       "测试需要升级的应用是否运行中，并获取应用的版本",
			Description: "通过访问应用统一的接口，正常返回说明运行中，否则不是",
			Action: func(c *cli.Context) {
				if alive := TestAlive(); alive == true {
					beego.Info("应用运行中")
					if len(G_UpdatedAppName) > 0 {
						DebugInfo(fmt.Sprintf("要升级的应用名称: %s(%s)", G_UpdatedAppName, G_currentUpdateInfo.Version) + GetFileLocation())
					}
				} else {
					beego.Info("应用没有运行")
				}
			},
		}, {
			Name:        "Kill",
			ShortName:   "kill",
			Usage:       "停止需要升级的应用",
			Description: "停止应用后，才能替换升级文件",
			Action: func(c *cli.Context) {
				if err := KillApp(); err != nil {
					DebugMust("关闭应用出错：" + err.Error() + GetFileLocation())
				} else {
					DebugMust("应用已成功关闭")
				}
			},
		},
	}
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Println("等待输入。。。")

			data, _, _ := reader.ReadLine()
			command := string(data)
			cliApp.Run(strings.Split(command, " "))
		}
	}()
	// app.Run(os.Args)
}
