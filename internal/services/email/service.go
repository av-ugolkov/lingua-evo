package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

//go:embed templ/*
var templ embed.FS

const (
	contentType = "Content-Type: text/html; charset=UTF-8\r\n\r\n"
)

type (
	authData struct {
		Code int
	}

	supportData struct {
		UserName string
		Msg      string
	}

	updatePsw struct {
		UserName string
		Code     int
	}
)

type Service struct {
	email config.Email
}

func NewService(email config.Email) *Service {
	return &Service{
		email: email,
	}
}

func (s *Service) SendAuthCode(toEmail string, code int) error {
	to := fmt.Sprintf("To: %s\r\n", toEmail)
	subject := "Subject: Create account on Lingua Evo\r\n"

	fs, err := template.ParseFS(templ, "templ/auth_code.html")
	if err != nil {
		return fmt.Errorf("email.Service.SendAuthCode: %v", err)
	}

	w := &bytes.Buffer{}
	err = fs.Execute(w, authData{
		Code: code,
	})
	if err != nil {
		return fmt.Errorf("email.Service.SendAuthCode: %v", err)
	}

	message := []byte(to + subject + contentType + w.String())

	authEmail := smtp.PlainAuth(runtime.EmptyString, s.email.Address, s.email.Password, s.email.Host)
	err = smtp.SendMail(s.email.AddrSvc(), authEmail, s.email.Address, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("email.Service.SendAuthCode: %v", err)
	}

	return nil
}

func (s *Service) SendEmailForSupport(toEmail string, userName, subject, msg string) error {
	to := fmt.Sprintf("To: %s\r\n", toEmail)
	subject = fmt.Sprintf("Subject: %s\r\n", subject)

	fs, err := template.ParseFS(templ, "templ/support.html")
	if err != nil {
		return fmt.Errorf("email.Service.SendEmailForSupport: %v", err)
	}

	if len(userName) == 0 {
		userName = "Sir/Madam"
	}

	w := &bytes.Buffer{}
	err = fs.Execute(w, supportData{
		UserName: userName,
		Msg:      msg,
	})
	if err != nil {
		return fmt.Errorf("email.Service.SendEmailForSupport: %v", err)
	}

	message := []byte(to + subject + contentType + w.String())

	authEmail := smtp.PlainAuth(runtime.EmptyString, s.email.Address, s.email.Password, s.email.Host)
	err = smtp.SendMail(s.email.AddrSvc(), authEmail, s.email.Address, []string{toEmail, s.email.Address}, message)
	if err != nil {
		return fmt.Errorf("email.Service.SendEmailForSupport: %v", err)
	}

	return nil
}

func (s *Service) SendEmailForUpdatePassword(toEmail, userName string, code int) error {
	to := fmt.Sprintf("To: %s\r\n", toEmail)
	subject := fmt.Sprintf("Subject: Update password\r\n")

	fs, err := template.ParseFS(templ, "templ/update_psw.html")
	if err != nil {
		return fmt.Errorf("email.Service.SendEmailForSupport: %v", err)
	}

	w := &bytes.Buffer{}
	err = fs.Execute(w, updatePsw{
		UserName: userName,
		Code:     code,
	})
	if err != nil {
		return fmt.Errorf("email.Service.SendEmailForSupport: %v", err)
	}

	message := []byte(to + subject + contentType + w.String())

	authEmail := smtp.PlainAuth(runtime.EmptyString, s.email.Address, s.email.Password, s.email.Host)
	err = smtp.SendMail(s.email.AddrSvc(), authEmail, s.email.Address, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("email.Service.SendEmailForSupport: %v", err)
	}

	return nil
}
