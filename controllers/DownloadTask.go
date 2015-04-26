package controllers

import (
	// "bufio"
	"github.com/Bluek404/downloader"
	// "github.com/astaxie/beego"
	// "github.com/codegangsta/cli"
	// "io"
	// "fmt"
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
