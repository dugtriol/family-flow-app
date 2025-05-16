package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"family-flow-app/config"
	"github.com/patrickmn/go-cache"
)

type EmailService struct {
	cache  *cache.Cache
	auth   smtp.Auth
	config config.Email
}

func NewEmailService(email config.Email) *EmailService {
	auth := smtp.PlainAuth("", email.FromEmail, email.Password, email.SMTP)

	c := cache.New(2*time.Minute, 5*time.Minute)
	return &EmailService{cache: c, auth: auth, config: email}
}

func (e *EmailService) SendCode(ctx context.Context, to []string) error {
	code := e.generateCode()

	subject := "Family Flow App - Код верификации"
	body := "Ваш код верификации: " + code
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

	// Сохраняем код в кэш
	e.cache.Set(to[0], code, cache.DefaultExpiration)

	return nil
}

func (e *EmailService) generateCode() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func (e *EmailService) CompareCode(ctx context.Context, email, code string) (bool, error) {

	val, found := e.cache.Get(email)
	if !found {
		return false, fmt.Errorf("code not found or expired")
	}

	if val != code {
		return false, nil
	}

	return true, nil
}

func (e *EmailService) GetAllKeys(ctx context.Context) ([]string, error) {

	keys := make([]string, 0)
	for k := range e.cache.Items() {
		keys = append(keys, k)
	}
	return keys, nil
}

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
