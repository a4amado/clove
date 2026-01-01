package emailTemplates

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"sync"
	"text/template"
)

//go:embed *.tmpl
var fs embed.FS

var once = &sync.Once{}
var templates *template.Template

// render_init initializes the package's template set by parsing all `*.tmpl` files from the embedded filesystem.
//
// It runs the initialization exactly once (safe for concurrent callers) and assigns the parsed templates to the package-level `templates` variable.
// The function will panic if template parsing fails.

type VerifyEmailTemplate struct {
	Token string
}

func (vet *VerifyEmailTemplate) formatEmailVerificationURL(token string) string {
	// TODO: make this url dynamic
	return fmt.Sprintf("https://clove.dev/api/v1/auth/verify-email?token=some-token%s", token)
}

func (vet *VerifyEmailTemplate) Render() (*string, error) {
	once.Do(func() {
		templates = template.Must(template.ParseFS(fs, "*.tmpl"))
	})

	buf := new(bytes.Buffer)
	err := templates.ExecuteTemplate(buf, "verify-email.tmpl", map[string]string{
		"verification_url": vet.formatEmailVerificationURL(vet.Token),
	})
	if err != nil {
		return nil, err
	}
	result := buf.String()
	return &result, nil
}
