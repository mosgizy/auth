package email

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v3"
)

type Service struct{
	client *resend.Client
	from string
}

func NewService(apiKey, fromEmail string) *Service{
	return &Service{
		client: resend.NewClient(apiKey),
		from: fromEmail,
	}
}

func (s *Service) SendResetPasswordEmail(ctx context.Context, to, resetLink string) error {
	params := &resend.SendEmailRequest{
		From: s.from,
		To: []string{to},
		Subject: "Reset Your Password",
		Html: fmt.Sprintf(`
			<h2>Reset Your Password</h2>
			<p>Click the link below to reset your password:</p>
			<p><a href="%s" style="color:#0066ff; font-size:18px;">Reset Password</a></p>
			<p>This link will expire in 15 minutes.</p>
			<p>If you didn't request this, please ignore this email.</p>
		`, resetLink),
	}

	_, err := s.client.Emails.SendWithContext(ctx,params)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}