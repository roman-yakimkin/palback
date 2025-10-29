package port

type EmailSender interface {
	SendVerificationEmail(toEmail, token string) error
	SendPasswordResetEmail(toEmail, token string) error
}
