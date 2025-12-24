package emailtemplates

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

func render_init() {
	once.Do(func() {
		templates = template.Must(template.ParseFS(fs, "*.tmpl"))
	})

}

type templateF interface {
	Render() string
}

type VerifyEmailTemplate struct {
	Token string
}

func (d *VerifyEmailTemplate) formatEmailverificationURL(token string) string {
	// TODO: make this url dynamic
	return fmt.Sprintf("https://clove.dev/api/v1/auth/verify-email?token=some-token%s", token)
}

func (d *VerifyEmailTemplate) Render() string {
	render_init()
	buf := new(bytes.Buffer)
	err := templates.ExecuteTemplate(buf, "verify-email.tmpl", map[string]string{
		"verification_url": d.formatEmailverificationURL(d.Token),
	})
	if err != nil {
		panic(err)
	}
	return buf.String()
}
