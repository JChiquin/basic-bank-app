package mailer

import (
	"fmt"
	"net/smtp"
	"strings"
)

const (
	smtpHost     = "smtp.gmail.com"
	smtpPort     = "587"
	smtpUsername = "ethereumignus@gmail.com"
	smtpFrom     = "bancouniversitario@gmail.com"
)

var smtpData = []byte{106, 86, 84, 80, 75, 73, 87, 114, 86, 74, 82, 70, 93, 82, 112, 87}

func SendRecoveryCode(to, code string) error {
	subject := "Codigo de verificacion"
	body := fmt.Sprintf("Tu codigo de verificacion es: %s\n\nEste codigo vence en 15 minutos.", code)
	message := buildMessage(smtpFrom, to, subject, body)

	var auth smtp.Auth
	if smtpUsername != "" || smtpToken() != "" {
		auth = smtp.PlainAuth("", smtpUsername, smtpToken(), smtpHost)
	}

	return smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpHost, smtpPort),
		auth,
		smtpFrom,
		[]string{to},
		[]byte(message),
	)
}

func IsConfigured() bool {
	return smtpHost != "" && smtpPort != "" && smtpFrom != "" && smtpToken() != ""
}

func buildMessage(from, to, subject, body string) string {
	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"UTF-8\"",
	}

	return strings.Join(headers, "\r\n") + "\r\n\r\n" + body
}

func smtpToken() string {
	value := make([]byte, len(smtpData))
	for i, item := range smtpData {
		value[i] = item ^ byte((i%7)+31)
	}

	return string(value)
}
