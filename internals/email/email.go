package email

import "github.com/mailjet/mailjet-apiv3-go/v4"

const publicKey = "3eb4240fd3a39737a465193758b8a60f"

var secretKey = "3136fc98959d7822de0bad6f3efa9ce7"

type Email struct {
	Client    *mailjet.Client
	ApiV      int
	FromEmail string
	FromName  string
}

type NewEmailClientOptions struct {
	FromEmail string
	FromName  string
	ApiV      int
}

func New(opts NewEmailClientOptions) Email {
	return Email{
		Client:    mailjet.NewMailjetClient(publicKey, secretKey),
		ApiV:      opts.ApiV,
		FromEmail: opts.FromEmail,
		FromName:  opts.FromName,
	}
}

func (e *Email) Auth() *AuthEmails {
	return &AuthEmails{
		email: e,
	}
}
