package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"net/smtp"
	// "os"
	// "os/exec"
)

var (
	DEFAULT_LOG_RECEIVER         = "errorlog@aliyun.com"
	DEFAULT_LOG_SENDER           = "errorlog@aliyun.com"
	DEFAULT_LOG_SENDER_PWD       = "1234567890qwertyuiop"
	DEFAULT_LOG_SENDER_SMTP_HOST = "smtp.aliyun.com"
)

//发送系统异常信息报告到指定邮箱
//
// appID 特定应用的标识
func EmailErrorLog(appID, body string, receiver []string) {
	if receiver == nil {
		receiver = []string{DEFAULT_LOG_RECEIVER}
		// receiver = []string{"ssor@qq.com"}
	}
	// Set up authentication information.
	auth := smtp.PlainAuth("", DEFAULT_LOG_RECEIVER, DEFAULT_LOG_SENDER_PWD, DEFAULT_LOG_SENDER_SMTP_HOST)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	// to := []string{"recipient@example.net"}
	msg := []byte(fmt.Sprintf("subject: 系统异常信息\r\n\r\n 来自[%s] \r\n %s", appID, body))
	go func() {
		err := smtp.SendMail(DEFAULT_LOG_SENDER_SMTP_HOST+":25", auth, DEFAULT_LOG_SENDER, receiver, msg)
		if err != nil {
			beego.Error(err)
		}
	}()

}

// ----------------------------------------------------------------------------------------
//接收异常信息的输出
type OutputTemp struct {
	temp string
}

func (this *OutputTemp) Write(p []byte) (n int, err error) {
	this.temp = this.temp + string(p)
	return len(p), nil
}
func (this *OutputTemp) String() string {
	return this.temp
}

// ----------------------------------------------------------------------------------------
