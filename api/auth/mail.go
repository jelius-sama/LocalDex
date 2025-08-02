package auth

import (
	"LocalDex/logger"
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

type SendMail struct {
	Host    string
	Port    string
	User    string
	Pass    string
	From    string
	To      string
	Subject string
	Body    string
}

// INFO: Fuck iCloud

// normalizeEmail trims spaces and returns a valid SMTP address like <addr>
func normalizeEmail(addr string) (string, error) {
	parsed, err := mail.ParseAddress(strings.TrimSpace(addr))
	if err != nil {
		return "", fmt.Errorf("invalid email format: %w", err)
	}
	return parsed.Address, nil
}

// sendMailSMTP587 handles the full STARTTLS+AUTH+SEND sequence on port 587.
func sendMailSMTP587(arg SendMail) error {
	logger.Debug(fmt.Sprintf(
		"SMTP Debug → Host: %s, Port: %s, User: %s, AuthHost: %s",
		arg.Host, arg.Port, arg.User, arg.Host,
	))

	addr := net.JoinHostPort(arg.Host, arg.Port)

	// Normalize addresses
	fromAddr, err := normalizeEmail(arg.From)
	if err != nil {
		return fmt.Errorf("from address error: %w", err)
	}
	toAddr, err := normalizeEmail(arg.To)
	if err != nil {
		return fmt.Errorf("to address error: %w", err)
	}

	// 1) Dial plain TCP
	c, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("dial smtp: %w", err)
	}
	defer c.Close()

	// 2) Say EHLO
	if err := c.Hello("localhost"); err != nil {
		return fmt.Errorf("helo: %w", err)
	}

	// 3) STARTTLS
	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsCfg := &tls.Config{ServerName: arg.Host}
		if err := c.StartTLS(tlsCfg); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	} else {
		return fmt.Errorf("server does not support STARTTLS")
	}

	// 4) AUTH
	auth := smtp.PlainAuth(arg.User, arg.User, arg.Pass, arg.Host)
	if err := c.Auth(auth); err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	// 5) MAIL FROM / RCPT TO (SMTP envelope)
	if err := c.Mail("<" + fromAddr + ">"); err != nil {
		return fmt.Errorf("mail from: %w", err)
	}
	if err := c.Rcpt("<" + toAddr + ">"); err != nil {
		return fmt.Errorf("rcpt to: %w", err)
	}

	// 6) DATA
	wc, err := c.Data()
	if err != nil {
		return fmt.Errorf("data: %w", err)
	}
	defer wc.Close()

	// 7) Write headers + body (use header display names if needed)
	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("From: %s\r\n", arg.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", arg.To))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", arg.Subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n")
	msg.WriteString(arg.Body)

	if _, err := wc.Write([]byte(msg.String())); err != nil {
		return fmt.Errorf("write msg: %w", err)
	}

	// 8) QUIT
	return c.Quit()
}
