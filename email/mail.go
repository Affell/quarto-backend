package email

import "html/template"

type Vars struct {
	Title      string
	Text       string
	ButtonText string
	ButtonURL  string
}

type Structure struct {
	To      string
	Subject string
	Vars    Vars
}

type Config struct {
	From            string
	User            string
	Password        string
	Host            string
	Port            string
	RecoverTemplate *template.Template
}

func New(to, subject, title, text, btnText, btnUrl string) Structure {
	return Structure{
		To:      to,
		Subject: subject,
		Vars: Vars{
			Title:      title,
			Text:       text,
			ButtonText: btnText,
			ButtonURL:  btnUrl,
		},
	}
}
