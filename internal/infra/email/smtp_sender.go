package email

import (
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"palback/internal/config"
)

type SMTPSender struct {
	config *config.Config
}

func NewSMTPSender(config *config.Config) *SMTPSender {
	return &SMTPSender{
		config: config,
	}
}

func (s *SMTPSender) SendVerificationEmail(toEmail, token string) error {
	verifyLink := fmt.Sprintf("%s/user/verify-email?token=%s", s.config.FrontendOrigin, token)

	// Простой HTML-шаблон
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head><meta charset="utf-8"></head>
	<body>
		<p>Здравствуйте!</p>
		<p>Пожалуйста, подтвердите вашу регистрацию, перейдя по ссылке:</p>
		<p><a href="{{.Link}}" style="display:inline-block;padding:10px 20px;background:#007bff;color:#fff;text-decoration:none;border-radius:4px;">Подтвердить регистрацию</a></p>
		<p>Ссылка действительна в течение 1 часа.</p>
		<hr>
		<p>С уважением,<br>Администрация проекта palomniki.su</p>
	</body>
	</html>
    `

	return s.sendMail(verifyLink, tmpl, "Подтверждение регистрации", toEmail)
}

func (s *SMTPSender) SendPasswordResetEmail(toEmail, token string) error {
	resetLink := fmt.Sprintf("%s/user/reset-password?token=%s", s.config.FrontendOrigin, token)

	// HTML-шаблон
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head><meta charset="utf-8"></head>
	<body>
		<p>Здравствуйте!</p>
		<p>Вы запросили сброс пароля. Если это были не вы — просто проигнорируйте это письмо.</p>
		<p>Чтобы установить новый пароль, перейдите по ссылке:</p>
		<p><a href="{{.Link}}" style="display:inline-block;padding:10px 20px;background:#007bff;color:#fff;text-decoration:none;border-radius:4px;">Сбросить пароль</a></p>
		<p>Ссылка действительна в течение 1 часа.</p>
		<hr>
		<p>С уважением,<br>Администрация проекта palomniki.su</p>
	</body>
	</html>
    `

	return s.sendMail(resetLink, tmpl, "Сброс пароля", toEmail)
}

func (s *SMTPSender) sendMail(link, tmpl, subject, toEmail string) error {
	t := template.Must(template.New("email").Parse(tmpl))
	var body strings.Builder
	if err := t.Execute(&body, struct{ Link string }{Link: link}); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\n"
	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		mime + "\r\n" +
		body.String())

	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	return smtp.SendMail(addr, auth, s.config.SMTPFrom, []string{toEmail}, msg)

}
