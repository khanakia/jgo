package auth

import (
	"github.com/khanakia/jgo/pkg/mail"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func SendRegisterEmail(user User, db *gorm.DB) (string, error) {

	if user.WelcomeEmailSent {
		return "", errors.New("Email was already sent")
	}

	data := map[string]interface{}{
		// "Name":         user.GetName(),
		"Email":        user.Email,
		"UID":          GetUserUID(user, db),
		"Password":     user.RandomPass,
		"OptinLink":    viper.GetString("api_url") + "/newsletter/optin?optin=true&uid=" + user.UID + "&redirect=",
		"MarketAppUrl": viper.GetString("market_app_url"),
	}

	files := []string{
		"user-register.html",
		"base-register.html",
	}

	body, err := mail.Compile(files, data)

	if err != nil {
		return "", err
	}

	msg := &mail.Message{
		To:      []string{user.Email},
		Subject: "Registration Successful",
		Body:    body,
	}
	mail.Send(msg)

	user.RandomPass = ""
	user.WelcomeEmailSent = true
	db.Save(&user)

	return body, nil
}

func SendForgotPasswordEmail(user User, secret string) (string, error) {

	// YTD
	data := map[string]interface{}{
		// "Name":  user.GetName(),
		"Token": secret,
		"Link":  viper.GetString("market_app_url") + "/auth/reset_password?email=" + user.Email + "&token=" + secret,
	}

	files := []string{
		"forgot-password.html",
		"base.html",
	}

	body, err := mail.Compile(files, data)

	if err != nil {
		return "", err
	}

	msg := &mail.Message{
		To:      []string{user.Email},
		Subject: "Reset Password",
		Body:    body,
	}
	mail.Send(msg)

	return body, nil
}
