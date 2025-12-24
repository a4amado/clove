package email

import (
	emailTemplates "clove/internals/email/email-templates"
	"context"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type AuthEmails struct {
	email *Email
}

type SendEmailVerificationToken struct {
	Token   string
	ToEmail string
	ToName  string
	Title   string
}

func (e *AuthEmails) SendEmailVerificaionToken(ctx context.Context, opts SendEmailVerificationToken) error {
	email := emailTemplates.VerifyEmailTemplate{
		Token: opts.Token,
	}
	emailHtml, err := email.Render()
	if err != nil {
		return err
	}
	_, err = e.email.Client.SendMailV31(&mailjet.MessagesV31{
		Info: []mailjet.InfoMessagesV31{
			{
				To: &mailjet.RecipientsV31{
					{
						Email: opts.ToEmail,
						Name:  opts.ToName,
					},
				},
				From: &mailjet.RecipientV31{
					Email: e.email.FromEmail,
					Name:  e.email.FromName,
				},
				Subject:  "Verify your email",
				HTMLPart: *emailHtml,
			},
		},
	}, mailjet.WithContext(ctx))
	if err != nil {
		return err
	}
	return nil
}
