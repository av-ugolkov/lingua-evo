package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/av-ugolkov/lingua-evo/internal/config"
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
		Msg string
	}
)

type Service struct {
	email  config.Email
	client *smtp.Client
}

func NewService(email config.Email) *Service {
	return &Service{
		email: email,
	}
}

func (s *Service) SendAuthCode(toEmail string, code int) error {
	to := fmt.Sprintf("To: %s\r\n", toEmail)
	subject := "Subject: Create account on Lingua Evo\r\n\r\n"

	fs, err := template.ParseFS(templ, "templ/auth_code.html")
	if err != nil {
		return fmt.Errorf("email.Service.SendAuthCode - parse template: %v", err)
	}

	w := &bytes.Buffer{}
	err = fs.Execute(w, authData{
		Code: code,
	})
	if err != nil {
		return fmt.Errorf("email.Service.SendAuthCode - execute template: %v", err)
	}

	message := []byte(to + subject + contentType + w.String())

	authEmail := smtp.PlainAuth("", s.email.Address, s.email.Password, s.email.Host)
	err = smtp.SendMail(s.email.AddrSvc(), authEmail, s.email.Address, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("email.Service.SendAuthCode - send mail: %v", err)
	}

	return nil
}

func (s *Service) SendEmailForSupport(toEmail string, params ...string) error {
	to := fmt.Sprintf("To: %s\r\n", toEmail)
	subject := fmt.Sprintf("Subject: Re: %s\r\n", params[1])

	fs, err := template.ParseFS(templ, "templ/support.html")
	if err != nil {
		return fmt.Errorf("email.Service.SendEmailForSupport - parse template: %v", err)
	}

	w := &bytes.Buffer{}
	err = fs.Execute(w, supportData{
		Msg: params[2],
	})
	if err != nil {
		return fmt.Errorf("email.Service.SendEmailForSupport - execute template: %v", err)
	}

	message := []byte(to + subject + contentType + w.String())

	authEmail := smtp.PlainAuth("", s.email.Address, s.email.Password, s.email.Host)
	err = smtp.SendMail(s.email.AddrSvc(), authEmail, s.email.Address, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("email.Service.SendEmailForSupport - send mail: %v", err)
	}

	return nil
}
