package controllers

import (
	// "fmt"
	"encoding/json"
	// "github.com/astaxie/beego"
)

type UpdateInfo struct {
	Version  string
	FileList FileChecksumList
}

func (this *UpdateInfo) Print() {
	DebugInfo("升级信息：版本号 " + this.Version + GetFileLocation())
	this.FileList.Print()
}
func (this *UpdateInfo) ToJson() ([]byte, error) {
	return json.Marshal(this)
}

func NewUpdateInfo(version string, list FileChecksumList) *UpdateInfo {
	return &UpdateInfo{
		Version:  version,
		FileList: list,
	}
}
