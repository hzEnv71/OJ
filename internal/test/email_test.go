package test

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
	"testing"
)

func TestSendEmail(t *testing.T) {
	e := &email.Email{
		To:      []string{"l2003hz@163.com"},
		From:    "Get <l2003hz@163.com>",
		Subject: "验证码已发送，请查收",
		Text:    []byte("Text Body is, of course, supported!"),
		HTML:    []byte("您的验证码：<b>" + "666666" + "</b>"),
		Headers: textproto.MIMEHeader{},
	}
	err := e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "l2003hz@163.com", "FPuvz3Byq7VV68eX", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		t.Fatal(err)
	}
}
