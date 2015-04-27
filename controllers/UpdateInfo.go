package controllers

import (
	// "fmt"
	"encoding/json"
	// "github.com/astaxie/beego"
)

type UpdateInfo struct {
	Version  string
	FileList FileChecksumList
	DirList  FileChecksumList //dir
}

func (this *UpdateInfo) ToJson() ([]byte, error) {
	return json.Marshal(this)
}

func (this *UpdateInfo) Print() {
	DebugSys("升级信息：版本号 " + this.Version)
	if len(this.DirList) > 0 {
		DebugSys("目录列表：")
		this.DirList.Print()
	} else {
		DebugSys("目录列表为空")
	}
	if len(this.FileList) > 0 {
		DebugSys("文件列表：")
		this.FileList.Print()
	} else {
		DebugSys("文件列表为空")
	}
}

func NewUpdateInfo(version string, list, dirList FileChecksumList) *UpdateInfo {
	return &UpdateInfo{
		Version:  version,
		FileList: list,
		DirList:  dirList,
	}
}
