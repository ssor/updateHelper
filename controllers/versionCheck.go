package controllers

import (
	// "bufio"
	// "fmt"
	// "github.com/Bluek404/downloader"
	// "github.com/astaxie/beego"
	// "github.com/astaxie/beego/config"
	// "github.com/codegangsta/cli"
	"io"
	// "os"
	// "strings"
	// "sync"
	"encoding/json"
	"errors"
	"net/http"
	// "os/exec"
	// "time"
	// "path"
	// "path/filepath"
)

func GetVersionInfo(url string) (*UpdateInfo, error) {
	bytes, err := RequestVersionInfo(url)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(bytes))
	var ui UpdateInfo
	if err := json.Unmarshal((bytes), &ui); err != nil {
		return nil, err
	}
	// ui.Print()
	return &ui, nil
}
func RequestVersionInfo(url string) ([]byte, error) {
	bytes := []byte{}
	resp, err := http.Get(url)
	// fmt.Println(resp.Status + "  " + url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()

	var buf = make([]byte, 1024)
	for {
		n, e := resp.Body.Read(buf)
		bytes = append(bytes, buf[:n]...)
		if e != nil {
			if e == io.EOF {
				// fmt.Println("数据读取完毕")
				// 数据已经下载完毕
				return bytes, nil
			}
			return nil, e
		}
	}
}
