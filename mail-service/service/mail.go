package service

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

type Mail struct {
}

func NewMailService() *Mail {
	return &Mail{}
}
func (service *Mail) SendEmail(to []string, content string) error {
	host := "smtp.gmail.com"
	auth := smtp.PlainAuth("", "duykhanh.forwork2108@gmail.com", "nndf jspt qnsb gtlk", host)
	title := "Report " + time.Now().Format("2006-01-02")
	header := make(map[string]string)
	header["From"] = host
	header["To"] = strings.Join(to, ",")
	header["Subject"] = title
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for i, j := range header {
		message += fmt.Sprintf("%s: %s\r\n", i, j)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(content))
	return smtp.SendMail("smtp.gmail.com:587", auth, "", to, []byte(message))
}
