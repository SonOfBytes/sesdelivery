package deliver

import (
	"io"
	netsmtp "net/smtp"
	"fmt"
)

type SMTP struct {
	client *netsmtp.Client
	server string
}

func NewSMTP(server string) (smtp *SMTP, err error) {
	return &SMTP{
		server: server,
	} , nil
}

func (s *SMTP) Send(sender string, recipients []string, body io.ReadCloser) (err error) {
	// Connect to the remote SMTP server.
	c, err := netsmtp.Dial(fmt.Sprintf("%s:25", s.server))
	if err != nil {
		return err
	}

	// Set the sender and recipient first
	if err := c.Mail(sender); err != nil {
		return err
	}
	for _, recipient := range recipients {
		if err := c.Rcpt(recipient); err != nil {
			return err
		}
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}
	_, err = io.Copy(wc, body)
	if err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		return err
	}
	return nil
}