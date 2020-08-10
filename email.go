package main

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
)

func GetSmtp(email string) string {
	if strings.Contains(email, "gmail") {
		return "gmail"
	}

	if strings.Contains(email, "outlook") {
		return "office365"
	}

	return ""
}

func Send(fromEmail, toEmail, smtpStr, pwd, name, msg string) error {

	from := fromEmail
	password := pwd

	to := []string{
		toEmail,
	}

	smtpHost := "smtp" + smtpStr + ".com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    name,
		Message: msg,
	})

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
