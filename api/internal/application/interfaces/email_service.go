package interfaces

type EmailService interface {
	SendInvitation(email, token, leagueName string) error
	SendGameNotification(emails []string, gameDate, location string, signupLink string) error
}
