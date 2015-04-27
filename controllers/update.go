package controllers

import (
	// "bufio"
	"fmt"
	"github.com/Bluek404/downloader"
	// "github.com/astaxie/beego"
	// "github.com/astaxie/beego/config"
	// "github.com/codegangsta/cli"
	"io"
	"os"
	"strings"
	// "sync"
	"encoding/json"
	// "errors"
	// "net/http"
	// "os/exec"
	"time"
	// "path"
	"path/filepath"
)

func CheckUpdate() {
	if ui, err := GetVersionInfo(G_versionUrl); err != nil {
		DebugSys(err.Error() + GetFileLocation())
		return
	} else {
		if ui.Version > G_currentUpdateInfo.Version { //版本有提升才有意义
			if G_downloadingUpdateInfo == nil {
				//没有正在下载的升级任务
				DebugInfo("系统有更新，准备升级" + GetFileLocation())
				clearUpdateResourceDir()
				StartDownloadTask(createDownloadTasks(ui))
			} else {
				if ui.Version > G_downloadingUpdateInfo.Version {
					//如果正在下载升级文件，并且版本号有提升
					//停止正在进行的下载，清空目录
					stopDownlingTasks()
					clearUpdateResourceDir()
					StartDownloadTask(createDownloadTasks(ui))
				}
			}
		} else {
			DebugInfo("系统无需升级" + GetFileLocation())
		}
	}
}
func clearUpdateResourceDir() {
	// os.Remove(G_UpdateResourcePath + G_versionInfoFile)
	// os.RemoveAll(G_UpdateResourcePath + "Bin/")
	os.RemoveAll(G_UpdateResourcePath)
}
func stopDownlingTasks() {
	G_downloadTasks.Abort()
}

//开始下载升级文件
func createDownloadTasks(ui *UpdateInfo) DownloadTaskList {
	//查看是否有新目录需要创建
	if err := PrepareUpdateFileDir(ui.DirList); err != nil {
		DebugMust(err.Error() + GetFileLocation())
		return nil
	}

	//比对当前版本与新版本的文件异同，只挑选与当前版本不同和当前版本缺少的文件下载
	needDownFiles := FileChecksumList{}
	for _, file := range ui.FileList {
		if G_currentUpdateInfo.FileList.Contains(file) == false {
			needDownFiles = needDownFiles.Add(file)
		}
	}

	downloadTasks := DownloadTaskList{}
	for _, fc := range needDownFiles {
		downloadTasks = append(downloadTasks, NewDownloadTask(G_UpdateResourcePath+"Bin/"+fc.Path, G_baseUrl+"Bin/"+fc.Path))
	}
	if len(downloadTasks) > 0 {
		G_downloadingUpdateInfo = ui //表示有升级任务下载
	}
	return downloadTasks
	// if err := StartDownloadTask(G_downloadTasks); err != nil {
	// 	DebugSys("下载升级文件出错" + GetFileLocation())
	// } else {
	// 	DebugInfo("升级文件准备完毕" + GetFileLocation())
	// }
}

//开始下载任务
func StartDownloadTask(tasks DownloadTaskList) error {
	if tasks == nil {
		return nil
	}
	taskCount := len(tasks)
	if taskCount <= 0 {
		DebugInfo("没有文件下载任务" + GetFileLocation())
		return nil
	}
	G_downloadTasks = tasks
	chanTask := make(chan DownloadTask, taskCount)
	for _, task := range tasks {
		go func(path, url string) {
			if fl, err := DownloadFromUrl(path, url, chanTask); err != nil {
				DebugMust(err.Error() + GetFileLocation())
			} else {
				task.Downloader = fl
			}
		}(task.Path, task.Url)
	}
	for i := 0; i < taskCount; i++ {
		downloadResult := <-chanTask
		tasks.SetStatus(downloadResult.Url, downloadResult.Status)
		DebugTrace(fmt.Sprintf("接收到下载反馈信息 %s url: %s", downloadResult.Status, downloadResult.Url) + GetFileLocation())
	}

	for _, task := range tasks {
		if task.Status == true {
			DebugInfo(fmt.Sprintf("%s 下载完毕", task.Url) + GetFileLocation())
		} else {
			DebugSys(fmt.Sprintf("%s 下载失败", task.Url) + GetFileLocation())
		}
	}
	completedTaskList := tasks.GetCompletdTaskList()
	if len(completedTaskList) < len(tasks) {
		DebugSys(fmt.Sprintf("下载升级文件出错，共有 %d 个下载任务，完成了 %d 个", len(tasks), len(completedTaskList)))
		G_downloadingUpdateInfo = nil //将其设为空，等下一次下载的时候会重复下载
		return G_errNotAllTaskCompleted
	}
	DebugSys("下载任务全部完成，升级文件准备完毕" + GetFileLocation())

	//将版本信息写入到文件中，表明这个版本升级文件准备完毕
	fileName := G_UpdateResourcePath + G_versionInfoFile
	if err := createVersionInfoFile(fileName, G_downloadingUpdateInfo); err != nil {
		return err
	}
	closeCheckUpdateIntervalMode()
	G_updateAppReady = true
	return nil
}
func UpdateApp() {
	G_updateAppReady = false
	KillApp()
	//将升级文件拷贝到系统目录中
	if err := copyUpdateFileToApp(); err != nil {
		DebugMust("升级中出错：" + err.Error() + GetFileLocation())
		openCheckUpdateIntervalMode()
		return
	} else {
		G_currentUpdateInfo = G_downloadingUpdateInfo
		G_downloadingUpdateInfo = nil
	}
	clearUpdateResourceDir()
	StartApp()
	openCheckUpdateIntervalMode()
	DebugMust("升级成功" + GetFileLocation())
}

//查看升级文件是否在新目录中，如果有先创建该目录
func PrepareUpdateFileDir(dirList FileChecksumList) error {
	// dirList := []string{}
	// for _, filePath := range dirList {
	// 	dirList = append(dirList, getDirDepList(filePath.Path, []string{})...)
	// }
	// dirList = removeDupDir(dirList, []string{})
	// DebugSys("需要创建的目录：" + GetFileLocation())
	// fmt.Println(dirList)
	for _, _dir := range dirList {
		dirpath := G_UpdateResourcePath + "Bin/" + _dir.Path
		DebugTrace(fmt.Sprintf("创建目录：%s", dirpath) + GetFileLocation())
		if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
func removeDupDir(dirList, noDupDirList []string) []string {
	if len(dirList) <= 0 {
		return noDupDirList
	} else {
		dir := dirList[0]
		dirList = dirList[1:]
		for _, _dir := range noDupDirList {
			if _dir == dir {
				return removeDupDir(dirList, noDupDirList)
			}
		}
		return removeDupDir(dirList, append(noDupDirList, dir))
	}
}

//根据传入的路径，将其按级分解成一个列表，如 a/b/c.js 分解成 a  a/b
func getDirDepList(fullPath string, paths []string) []string {
	// fmt.Println("fullPath: " + fullPath)
	dir, _ := filepath.Split(fullPath)
	// fmt.Println("dir: " + dir)
	if dir == "/" || dir == "." || dir == "./" || len(dir) <= 0 {
		return paths
	} else {
		paths = append([]string{dir}, paths...)
		// fmt.Println(paths)
		return getDirDepList(filepath.Dir(dir), paths)
	}
}

func copyUpdateFileToApp() error {
	walkFn := func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(info.Name(), ".") == true {
			return nil
		}
		if info.IsDir() == true { //在应用目录内查看是否已创建该目录
			dirPath := strings.Replace(fullPath, "UpdateResource", "App", 1)
			DebugInfo(fmt.Sprintf("需要创建目录：%s", dirPath) + GetFileLocation())
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				DebugMust(fmt.Sprintf("创建目录 %s 失败：%s", dirPath, err.Error()) + GetFileLocation())
				return err
			} else {
				DebugInfo(fmt.Sprintf("创建目录 %s 成功", dirPath) + GetFileLocation())
			}
		} else { //将文件拷贝到应用对应目录内
			destFilePath := strings.Replace(fullPath, "UpdateResource", "App", 1)
			DebugInfo(fmt.Sprintf("从 %s 向 %s 拷贝文件", fullPath, destFilePath) + GetFileLocation())
			if Exist(destFilePath) == true {
				if err := os.Remove(destFilePath); err != nil {
					return err
				}
			}
			// DebugInfo(fmt.Sprintf("需要复制文件：%s", destFilePath) + GetFileLocation())
			if err := CopyFile(destFilePath, fullPath); err != nil {
				DebugMust(fmt.Sprintf("从 %s 向 %s 拷贝文件出错：%s", fullPath, destFilePath, err.Error()) + GetFileLocation())
				return err
			} else {
				DebugInfo(fmt.Sprintf("从 %s 向 %s 拷贝文件成功", fullPath, destFilePath) + GetFileLocation())
			}
		}
		return nil
	}
	if err := filepath.Walk(G_UpdateResourcePath, walkFn); err != nil {
		return err
	}
	return nil
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
func CopyFile(dstName, srcName string) error {
	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}
func createVersionInfoFile(fileName string, versionInfo *UpdateInfo) error {
	writeFile := func(bytes []byte) error {
		if fd, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755); err != nil {
			return err
		} else {
			if _, err := fd.Write(bytes); err != nil {
				return err
			} else {
				fd.Close()
			}
		}
		return nil
	}
	if bytes, err := json.Marshal(versionInfo); err != nil {
		return err
	} else {
		if err := writeFile(bytes); err != nil {
			DebugMust("创建版本信息文件失败：" + err.Error() + GetFileLocation())
			return err
		}
	}
	return nil
}
func openCheckUpdateIntervalMode() {
	G_CheckUpdateIntervalMode = true //
}
func closeCheckUpdateIntervalMode() {
	G_CheckUpdateIntervalMode = false //
}

func DownloadFromUrl(filePath, url string, chanDownloadCount chan DownloadTask) (*downloader.FileDl, error) {
	// fileTempPath := "./tmp/"
	//如果路径中包含文件夹，需要首先建立该文件夹
	// fileName := path.Base(url)
	// file, err := os.OpenFile(fileTempPath+fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		DebugSys(err.Error() + GetFileLocation())
		return nil, err
	}
	fileDl, err := downloader.NewFileDl(url, file, -1)
	if err != nil {
		DebugSys(fmt.Sprintf("下载 [%s] 出错：%s  url: %s", filePath, err.Error(), url) + GetFileLocation())
		chanDownloadCount <- DownloadTask{
			Url:    url,
			Status: false,
		}
		// os.Remove(fileTempPath + fileName)
		os.Remove(filePath)
		return nil, err
	}
	var chanExit = make(chan bool)
	var chanProgress = make(chan bool)
	var chanAbort = make(chan bool)
	fileDl.OnStart(func() {
		// fmt.Println("开始下载")
		for {
			select {
			case <-chanExit:
				status := fileDl.GetStatus()
				// fmt.Println(fmt.Sprintf(format, status.Downloaded, fileDl.Size, h, 0, "[FINISH]"))
				DebugInfo(fmt.Sprintf("[%s] 下载完成，共 %d 字节", filePath, status.Downloaded) + GetFileLocation())
				// DebugTrace("关闭文件"+GetFileLocation())
				file.Close()
				chanDownloadCount <- DownloadTask{
					Url:    url,
					Status: true,
				}
				return
			case <-chanAbort:
				i := 0
				for {
					if err := file.Close(); err == nil {
						DebugInfo("下载取消成功，关闭了文件 [" + filePath + "]" + GetFileLocation())
						break
					}
					time.Sleep(time.Second * 1)
					i++
					if i > 3 {
						DebugMust("下载取消失败，无法关闭文件 [" + filePath + "]" + GetFileLocation())
						break
					}
				}
				return
			case <-chanProgress:
				// format := "\033[2K\r%v/%v [%s] (当前速度： %v byte/s) %v"
				// status := fileDl.GetStatus()
				// var i = float64(status.Downloaded) / float64(fileDl.Size) * 50
				// if i < 0 {
				// 	i = 0
				// }
				// h := strings.Repeat("=", int(i)) + strings.Repeat(" ", 50-int(i))

				// fmt.Println(fmt.Sprintf(format, status.Downloaded, fileDl.Size, h, status.Speeds, "[DOWNLOADING]"))
			}
		}
	})
	fileDl.OnAbort(func() {
		chanAbort <- true
	})
	fileDl.OnProgress(func() {
		chanProgress <- true
	})
	fileDl.OnFinish(func() {
		chanExit <- true
	})

	fileDl.OnError(func(errCode int, err error) {
		fmt.Println(errCode, err)
		chanDownloadCount <- DownloadTask{
			Url:    url,
			Status: false,
		}
	})

	// fmt.Printf("%+v\n", fileDl)
	DebugInfo(fmt.Sprintf("开始下载 url: %s", url) + GetFileLocation())
	fileDl.Start()
	return fileDl, nil
}
