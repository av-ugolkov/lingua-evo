package support

import (
	"context"
	"fmt"
)

type (
	emailSvc interface {
		SendEmailForSupport(email string, params ...string) error
	}
)

type Service struct {
	emailSvc emailSvc
}

func NewService(emailSvc emailSvc) *Service {
	return &Service{
		emailSvc: emailSvc,
	}
}

func (s *Service) SendRequest(ctx context.Context, request SupportRequest) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("support.Service.SendRequest: %w", err)
		}
	}()
	err = s.emailSvc.SendEmailForSupport(request.Email, request.Name, request.Type, request.Message)
	if err != nil {
		return err
	}
	return nil
}
