package main

import (
	"github.com/astaxie/beego"
	_ "updateHelper/routers"
)

/*
	目录结构
	updateHelper
	App/
		Version.md（用来记录系统版本信息）
		Bin/
	UpdateResource/
		Version.md（用来标识是否升级文件下载完整）
		Bin/

	----------------------------------------------------------------------------------------------
	#启动App
	1. 检查可以启动的条件
		* App目录是否有版本信息文件
		* 读取App目录的版本信息文件，没有则使用 0.0，文件列表为空
	2. 如果满足启动条件，启动App

	------

	#定时更新升级信息（默认开启）
	1. 访问服务器检查升级信息，
		* 如果版本相同则略过，等待下一次定时
		* 如果高于当前版本，则比对缺失和不同的文件，清空Bin目录， 开始下载需要升级的文件
	2. 下载文件，完成后创建版本信息文件
	3. 如果下一次定时更新发生时，上一次更新没有结束，则检查两次更新的版本号是否一致，
		* 如果两次更新的版本号一致，则本次更新忽略
		* 如果两次更新的版本号不一致，停止上次更新，清空下载的文件，重新开始下载更新文件

	------

	#升级（默认不开启）
	1. 升级文件准备完毕后，开启定时升级模式，关闭定时检查服务器升级信息模式，设置为应用可升级状态，每隔一段时间向应用发送可升级通知
	2. 系统告知helper可以升级
	3. helper关闭系统，拷贝和替换文件，并删除已经拷贝完的文件
	4. 启动系统，
	5. 开启定时检查服务器升级信息设置

	----------------------------------------------------------------------------------------------
	开发路线：
	* 可以自动升级
	* 支持与应用配合升级


	//////////////////////////////////////////////////////////////////////////////////////////////


*/

func main() {

	// beego.SetStaticPath("/images", "static/images")
	// beego.SetStaticPath("/bootstrap", "static/bootstrap")
	// beego.SetStaticPath("/dataTable", "static/dataTable")
	// beego.SetStaticPath("/javascripts", "static/javascripts")
	// // beego.SetStaticPath("/layer", "static/layer")
	// beego.SetStaticPath("/stylesheets", "static/stylesheets")

	beego.Run()
}
