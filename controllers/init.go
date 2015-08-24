package controllers

import (
	"bufio"
	"fmt"
	"io/ioutil"
	// "github.com/Bluek404/downloader"
	"github.com/astaxie/beego"
	// "github.com/astaxie/beego/config"
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
	"github.com/BurntSushi/toml"
	"github.com/c4milo/unzipit"
	"github.com/ungerik/go-dry"
	// "runtime"
)

//系统配置项
type Config struct {
	AppID               string
	UpdatedAppPort      string
	UpdatedAppName      string
	UpdateServerBaseURL string
	UpdateCheckInterval int
}

func (this *Config) ListName() string {
	return "系统配置列表："
}
func (this *Config) InfoList() []string {
	list := []string{
		fmt.Sprintf("应用标识：%s", this.AppID),
		fmt.Sprintf("应用名称：%s", this.UpdatedAppName),
		fmt.Sprintf("应用升级端口：%s", this.UpdatedAppPort),
		fmt.Sprintf("应用升级资源URL：%s", this.UpdateServerBaseURL),
		fmt.Sprintf("应用升级检查时间间隔：%d 秒", this.UpdateCheckInterval),
	}
	return list
}

var (
	G_conf Config
)

func init() {
	// copyUpdateFileToApp()
	unzipApp("App.zip")
	if err := initConfig(); err != nil {
		return
	}
	go initCli()

	tryToStartApp()
	time.Sleep(time.Second * 2)

	// CheckUpdate()
	// return

	go startIntervalCheckUpdateInfoFromServer()
	go startIntervalNotifyAppUpdate()
}
func startIntervalNotifyAppUpdate() {
	c := time.Tick(15 * time.Second)
	// c := time.Tick(1 * time.Minute)
	for range c {
		if G_updateAppReady == true {
			go func() {
				DebugInfoF("提示App可以升级了")
				resp, err := http.Get("http://localhost:" + G_conf.UpdatedAppPort + "/Update")
				if err != nil {
					DebugSysF("App无法接收升级提示: %s", err.Error())
				}
				if resp.StatusCode != 200 {
					DebugSysF("App无法接收升级提示: %s", resp.Status)
				}
			}()
		} else {
			DebugInfoF("没有升级信息可以通知App")
		}
	}
}
func startIntervalCheckUpdateInfoFromServer() {
	c := time.Tick(time.Duration(G_conf.UpdateCheckInterval) * time.Second)
	// c := time.Tick(1 * time.Minute)
	for range c {
		if G_CheckUpdateIntervalMode == true {
			go func() {
				DebugInfoF("定时检查升级模式启用")
				CheckUpdate()
			}()
		} else {
			DebugInfoF("定时检查升级模式关闭")
		}
	}
}
func unzipApp(zipfileName string) {
	// dir := strings.Trim(zipfileName, path.Ext(zipfileName))
	// dir := "." + string(os.PathSeparator)
	dir := "./App"
	//处理压缩的数据，解压到当前同名的目录下
	if dry.FileExists(zipfileName) == true {
		if file, err := os.Open(zipfileName); err != nil {
			DebugMustF("打开数据库压缩文件出错：%s", err.Error())
			return
		} else {
			defer func() {
				if file != nil {
					file.Close()
				}
			}()
			DebugTraceF("解压缩App文件...")
			if _, errUnzip := unzipit.Unpack(file, dir); errUnzip != nil {
				DebugMustF("解压缩zip文件出错：%s", errUnzip.Error())
				return
			} else {
				file.Close()
				file = nil
				DebugTraceF("正在清理压缩文件 %s...", zipfileName)
				if err := os.Remove(zipfileName); err != nil {
					DebugMustF("清理压缩文件时出错：%s", err.Error())
					return
				}
				DebugTraceF("压缩文件 %s 清理完毕", zipfileName)
			}
		}
	}
}
func initConfig() error {
	confFile := "conf/sys.toml"
	if confData, err := ioutil.ReadFile(confFile); err != nil {
		DebugMustF("系统配置出错：%s", err.Error())
		return err
	} else {
		if _, err := toml.Decode(string(confData), &G_conf); err != nil {
			DebugMustF("系统配置出错：%s", err.Error())
			return err
		}
		DebugPrintList_Info(&G_conf)
	}

	return nil
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
				GetVersionInfo(G_conf.UpdateServerBaseURL + G_versionInfoFile)
			},
		}, {
			Name:        "StartApp",
			ShortName:   "sa",
			Usage:       "启动要升级的应用",
			Description: "要启动的应用位于 " + G_appBasePath + " 目录下",
			Action: func(c *cli.Context) {
				go TickForStartApp(3)

			},
		}, {
			Name:        "TestAlive",
			ShortName:   "ta",
			Usage:       "测试需要升级的应用是否运行中，并获取应用的版本",
			Description: "通过访问应用统一的接口，正常返回说明运行中，否则不是",
			Action: func(c *cli.Context) {
				if alive := TestAlive(); alive == true {
					beego.Info("应用运行中")
					if len(G_conf.UpdatedAppName) > 0 {
						DebugInfoF("要升级的应用名称: %s(%s)", G_conf.UpdatedAppName, G_currentUpdateInfo.Version)
					}
				} else {
					DebugInfoF("应用没有运行")
				}
			},
		}, {
			Name:        "Kill",
			ShortName:   "kill",
			Usage:       "停止需要升级的应用",
			Description: "停止应用后，才能替换升级文件",
			Action: func(c *cli.Context) {
				if err := KillApp(); err != nil {
					DebugMustF("关闭应用出错：%s", err.Error())
				} else {
					DebugMustF("应用已成功关闭")
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
