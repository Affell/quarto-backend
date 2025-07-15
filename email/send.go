package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html"
	"net"
	"net/mail"
	"net/smtp"
)

func (email Structure) Send(config Config) (err error) {

	from := mail.Address{
		Name:    config.From,
		Address: config.User,
	}

	to := mail.Address{
		Name:    "",
		Address: email.To,
	}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = html.EscapeString(email.Subject)
	headers["MIME-version"] = "1.0;"
	headers["Content-Type"] = `text/html; charset="utf-8"`

	var message string

	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	buff := new(bytes.Buffer)
	err = config.RecoverTemplate.ExecuteTemplate(buff, "email", email.Vars)
	if err != nil {
		err = fmt.Errorf("during executing template for email : %v", err)
		return
	}

	message += "\r\n" + buff.String()

	auth := login(config.User, config.Password)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         config.Host,
	}

	conn, err := net.Dial("tcp", net.JoinHostPort(config.Host, config.Port))
	if err != nil {
		err = fmt.Errorf("during net dial : %v", err)
		return
	}

	defer conn.Close()

	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		err = fmt.Errorf("during smtp.NewClient : %v", err)
		return
	}

	defer client.Quit()

	if err = client.StartTLS(tlsconfig); err != nil {
		err = fmt.Errorf("during client.StartTLS : %v", err)
		return
	}

	if err = client.Auth(auth); err != nil {
		err = fmt.Errorf("during client.Auth : %v", err)
		return
	}

	if err = client.Mail(from.Address); err != nil {
		err = fmt.Errorf("during client.Mail : %v", err)
		return
	}

	if err = client.Rcpt(to.Address); err != nil {
		err = fmt.Errorf("during client.Rcpt : %v", err)
		return
	}

	w, err := client.Data()
	if err != nil {
		err = fmt.Errorf("during client.Data : %v", err)
		return
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		err = fmt.Errorf("writing message in reader : %v", err)
		return
	}

	err = w.Close()
	if err != nil {
		err = fmt.Errorf("closing connection to the email server : %v", err)
		return
	}

	return
}
