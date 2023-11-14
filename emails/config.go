package emails

import (
	"apis-sati/utils"
	"crypto/tls"
	"github.com/go-gomail/gomail"
)

func MountEmail(body, title, email string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", utils.FROM)
	m.SetHeader("To", utils.CheckToSend(email))
	m.SetHeader("Subject", title)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(utils.SMTP, utils.PORT, utils.USERNAME, utils.PASSWORD)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)
	if err != nil {
		return err
	}
	return err
}
