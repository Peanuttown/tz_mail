package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Peanuttown/tz_mail/config"
	"github.com/jordan-wright/email"
	"io/ioutil"
	"net/smtp"
	"os"
	"path"
	"strings"
)

func main() {
	var to string
	var content string
	var subject string
	var attachfile string
	var help bool
	flag.StringVar(&to, "to", "", "send to ")
	flag.StringVar(&content, "c", "", "context")
	flag.StringVar(&subject, "s", "", "subject")
	flag.StringVar(&attachfile, "f", "", "attach file")
	flag.BoolVar(&help, "h", false, "help")
	flag.Parse()
	c := config.Config{}
	if help {
		flag.PrintDefaults()
		bytes, _ := json.Marshal(c)
		fmt.Println("config demo:\n", string(bytes))
		return
	}

	bytes, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".tz_mail"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	values := strings.Split(c.SMTPServerAddr, ":")
	if len(values) != 2 {
		fmt.Fprintln(os.Stderr, "服务器地址格式错误: ", c.SMTPServerAddr)
		os.Exit(1)
	}

	e := email.NewEmail()
	//e.From = "tzzNotify@163.com"
	e.From = c.User
	e.To = []string{to}
	e.Subject = subject
	if len(content) == 0 {
		content = "empty"
	}
	e.Text = []byte(content)
	if len(attachfile) > 0 {
		file, err := os.Open(attachfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		_, err = e.Attach(file, path.Base(attachfile), "")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	//e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	//e.Send("smtp.163.com:25", smtp.PlainAuth("", "tzzNotify@163.com", "AVIGHTZYNSINDXZX", "smtp.163.com"))
	err = e.Send(c.SMTPServerAddr, smtp.PlainAuth("", c.User, c.Passwd, values[0]))
	if err != nil {
		fmt.Fprintln(os.Stderr, "发送失败: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("发送成功")
}
