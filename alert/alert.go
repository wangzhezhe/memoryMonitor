package alert

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strconv"
	"strings"
)

//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//refer to http://www.oschina.net/code/snippet_166520_34694
//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}

func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func MemoryAlert(nodeIP string, containerID string, memoryPercentage int) error {

	host := "smtp.163.com:25"
	user := "18767169274@163.com"
	password := "cformalert123456"
	mailList := []string{"zhewang@daocloud.io", "zg.zhu@daocloud.io", "davidz@cform.io"}
	alertInfo := "Cform production alert!!! " + "nodeIP:" + nodeIP + " ContainerID: " + containerID + " MemoryPercentage: " + strconv.Itoa(memoryPercentage) + "%"
	subject := "Alert from daocloud"
	body := `
    <html>
    <body>
    <h3>
    ` + alertInfo + `
    </h3>
    </body>
    </html>
    `
	fmt.Println("alert")
	for _, value := range mailList {
		err := SendMail(user, password, host, value, subject, body, "html")
		if err != nil {
			return err
		} else {
			fmt.Printf("send mail to %s successfully!\n", value)
		}
	}

	return nil
}
