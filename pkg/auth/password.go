package auth

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/khanakia/jgo/pkg/cache"
	"github.com/khanakia/jgo/pkg/util"
	"gorm.io/gorm"
)

type ForgotPasswordRequest struct {
	UserName string `json:"userName"`
}

// Validate ...
func (a ForgotPasswordRequest) Validate() ([]interface{}, error) {
	errorValidate := validation.ValidateStruct(&a,
		validation.Field(&a.UserName, validation.Required),
	)
	data := util.ErroObjToArray(errorValidate)
	return data, errorValidate
}

func ForgotPassword(req ForgotPasswordRequest, cache cache.Store, db *gorm.DB) error {

	user := FindByEmail(req.UserName, db)
	if user == nil {
		return errors.New("User not found.")
	}

	secret := uuid.New().String()
	token := "fp:" + secret

	cache.Put(token, user.ID, 1000)

	SendForgotPasswordEmail(*user, secret)

	return nil
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// Validate ...
func (a ResetPasswordRequest) Validate() ([]interface{}, error) {
	errorValidate := validation.ValidateStruct(&a,
		validation.Field(&a.Token, validation.Required),
		validation.Field(&a.Password, validation.Required),
	)
	data := util.ErroObjToArray(errorValidate)
	return data, errorValidate
}

func (a *Auth) ResetPassword(req ResetPasswordRequest, cache cache.Store) error {
	cacheKey := "fp:" + req.Token
	userID := cache.Get(cacheKey)

	if userID == nil {
		return errors.New("Wrong token supplied")
	}

	var user User
	res := a.Dbc.Db.First(&user, userID)
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return errors.New("User does not exists")
	}

	user.Password = GeneratePassword(req.Password)
	err := a.Dbc.Db.Save(&user).Error

	if err != nil {
		return errors.New("Server error")
	}

	cache.Del(cacheKey)

	return nil
}
