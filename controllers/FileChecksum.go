package controllers

import (
	"fmt"
	// "github.com/astaxie/beego"
)

type FileChecksum struct {
	Path     string
	Checksum string
}

func (this *FileChecksum) Print() {
	DebugInfo(fmt.Sprintf("Path: %s  check: %s", this.Path, this.Checksum))
}
func NewFileChecksum(path, checksum string) *FileChecksum {
	return &FileChecksum{
		Path:     path,
		Checksum: checksum,
	}
}

type FileChecksumList []*FileChecksum

func (this FileChecksumList) Add(fc *FileChecksum) FileChecksumList {
	return append(this, fc)
}
func (this FileChecksumList) Contains(fc *FileChecksum) bool {
	for _, _fc := range this {
		if _fc.Path == fc.Path && _fc.Checksum == fc.Checksum {
			return true
		}
	}
	return false
}

func (this FileChecksumList) Print() {
	for _, fc := range this {
		fc.Print()
	}
}
