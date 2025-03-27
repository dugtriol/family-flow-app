package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"family-flow-app/config"
	"family-flow-app/pkg/redis"
)

type EmailService struct {
	rd     *redis.Redis
	auth   smtp.Auth
	config config.Email
}

func NewEmailService(rd *redis.Redis, email config.Email) *EmailService {
	auth := smtp.PlainAuth("", email.FromEmail, email.Password, email.SMTP)
	return &EmailService{rd: rd, auth: auth, config: email}
}

func (e *EmailService) SendCode(ctx context.Context, to []string) error {
	code := e.generateCode()

	subject := "Family Flow App - Код верификации"
	body := "Ваш код верификации: " + code
	//message := "Subject: " + subject + "\r\n" + body
	message := "Subject: " + subject + "\r\n" + "Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n" + body

	if err := smtp.SendMail(
		e.config.Addr,
		e.auth,
		e.config.FromEmail,
		to,
		[]byte(message),
	); err != nil {
		return err
	}

	statusCmd := e.rd.Set(ctx, to[0], code, 2*time.Minute)
	if err := statusCmd.Err(); err != nil {
		return err
	}

	return nil
}

func (e *EmailService) generateCode() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func (e *EmailService) CompareCode(ctx context.Context, email, code string) (bool, error) {
	statusCmd := e.rd.Get(ctx, email)
	if err := statusCmd.Err(); err != nil {
		return false, ErrCode
	}

	if statusCmd.Val() != code {
		return false, nil
	}

	return true, nil
}

func (e *EmailService) GetAllKeys(ctx context.Context) ([]string, error) {
	keysCmd := e.rd.Keys(ctx, "*")
	if err := keysCmd.Err(); err != nil {
		return nil, err
	}
	return keysCmd.Val(), nil
}

// отправить пригласительное письмо
func (e *EmailService) SendInvite(ctx context.Context, invite InputSendInvite) error {
	subject := "Family Flow App - Приглашение присоединиться к семье"
	body := fmt.Sprintf(
		"Пользователь %s (%s) приглашает присоединиться к семье %s. Перейдите по ссылке для регистрации.",
		invite.FromName,
		invite.From,
		invite.FamilyName,
	)
	message := "Subject: " + subject + "\r\n" + "Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n" + body

	if err := smtp.SendMail(
		e.config.Addr,
		e.auth,
		e.config.FromEmail,
		invite.To,
		[]byte(message),
	); err != nil {
		return err
	}

	return nil
}
