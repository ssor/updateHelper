package controllers

import (
	// "bufio"
	"github.com/Bluek404/downloader"
	// "github.com/astaxie/beego"
	// "github.com/codegangsta/cli"
	// "io"
	"fmt"
	// "os"
	// "strings"
	// "sync"
	// "time"
	// "errors"
	// "path"
)

type DownloadTask struct {
	Path       string
	Url        string
	Status     bool
	Downloader *downloader.FileDl
}

// func (this *DownloadTask) Print() {
// 	// DebugInfo(fmt.Sprintf("Path: %s  Url: %s", this.Path, this.Url) + GetFileLocation())
// 	DebugInfoF("Path: %s  ", this.Path)
// }
func (this *DownloadTask) String() string {
	return fmt.Sprintf("Path: %s", this.Path)
}
func (this *DownloadTask) Abort() {
	if this.Downloader != nil {
		this.Downloader.Abort()
	}
}
func NewDownloadTask(path, url string) *DownloadTask {
	return &DownloadTask{
		Path:   path,
		Url:    url,
		Status: false,
	}
}

type DownloadTaskList []*DownloadTask

func (this DownloadTaskList) ListName() string {
	return fmt.Sprintf("下载任务列表(%d)", len(this))
}
func (this DownloadTaskList) InfoList() []string {
	list := []string{}
	for _, task := range this {
		list = append(list, task.String())
	}
	return list
}

// func (this DownloadTaskList) Print() {
// 	if len(this) > 0 {
// 		DebugSysF("下载任务列表中有 %d 个任务", len(this))

// 		for _, task := range this {
// 			task.Print()
// 		}
// 	} else {
// 		DebugSysF("下载任务列表中没有任务")
// 	}
// }

func (this DownloadTaskList) SetStatus(url string, status bool) {
	for _, task := range this {
		if task.Url == url {
			task.Status = status
		}
	}
}
func (this DownloadTaskList) GetCompletdTaskList() DownloadTaskList {
	list := DownloadTaskList{}
	for _, task := range this {
		if task.Status == true {
			list = append(list, task)
		}
	}
	return list
}
func (this DownloadTaskList) Abort() {
	for _, task := range this {
		task.Abort()
	}
}
