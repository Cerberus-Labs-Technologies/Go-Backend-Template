package mail

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/wneessen/go-mail"
	"html/template"
	"necross.it/backend/auth/user"
	"necross.it/backend/util"
	"os"
)

//go:embed views/verify.html
var verifyView string
var verifyTemplate = template.Must(template.New("verify").Parse(verifyView))

//go:embed views/deleteAccount.html
var deleteAccountView string
var deleteAccountTemplate = template.Must(template.New("deleteAccount").Parse(deleteAccountView))

//go:embed views/forgotPassword.html
var forgotPasswordView string
var forgotPasswordTemplate = template.Must(template.New("forgotPassword").Parse(forgotPasswordView))

func SendDeleteAccountMail(user user.User) error {
	var executedTemplate bytes.Buffer
	err := deleteAccountTemplate.Execute(&executedTemplate, &user)
	go SendMail(user.Name, user.Email, "PricelessKeys - Account", executedTemplate.String())
	return err
}

func SendForgotPasswordMail(user user.User, forgotPassword ForgotPasswordMail) error {
	var executedTemplate bytes.Buffer
	err := forgotPasswordTemplate.Execute(&executedTemplate, &forgotPassword)
	go SendMail(user.Name, user.Email, "PricelessKeys - Forgot Password", executedTemplate.String())
	return err
}

// TODO: Send mail placeholder
func SendVerifyMail(user user.User, verify MailVerify) error {
	var executedTemplate bytes.Buffer
	verify.Name = user.Name
	err := verifyTemplate.Execute(&executedTemplate, &verify)
	go SendMail(user.Name, user.Email, "PricelessKeys - Verify", executedTemplate.String())
	return err
}

func SendMail(receiverName string, receiverMail string, subject string, body string) error {
	// Create a new mail message
	m := mail.NewMsg()

	if err := m.FromFormat("PricelessKeys", "noreply@pricelesskeys.com"); err != nil {
		fmt.Printf("failed to set FROM address: %s\n", err)
	}
	receiverString := receiverName + " <" + receiverMail + ">"
	if err := m.To(receiverString); err != nil {
		fmt.Printf("failed to set TO address: %s\n", err)
		os.Exit(1)
	}

	m.Subject(subject)
	m.SetDate()
	m.SetMessageID()
	m.SetBulk()
	m.SetBodyString(mail.TypeTextHTML, body)

	host := os.Getenv("MAIL_HOST")
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	c, err := mail.NewClient(host,
		mail.WithSMTPAuth(mail.SMTPAuthLogin), mail.WithUsername(username),
		mail.WithPassword(password), mail.WithTLSPolicy(mail.TLSMandatory), mail.WithPort(587))
	if err != nil {
		fmt.Printf("failed to create client: %s\n", err)
	}

	if err := c.DialAndSend(m); err != nil {
		fmt.Printf("failed to send mail: %s\n", err)
	}

	return err
}

type ForgotPasswordMail struct {
	Link string
}

type MailVerify struct {
	ID          int            `json:"id" db:"ID"`
	UserID      int            `json:"userId" db:"userId"`
	Name        string         `json:"name"`
	VerifyToken string         `json:"verifyToken" db:"verifyToken"`
	CreatedAt   util.TimeStamp `json:"createdAt" db:"createdAt"`
}

func (m MailVerify) GetVerifyLink() string {
	return "https://pricelesskeys.com/verify/" + m.VerifyToken
}
