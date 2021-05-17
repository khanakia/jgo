package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Claims ...
type Claims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	UID   string `json:"uid"`
	jwt.StandardClaims
}

type LoginArgs struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// Login ...
func Login(req LoginArgs, db *gorm.DB) (map[string]interface{}, error) {
	user := FindByEmail(req.UserName, db)
	if user == nil {
		return nil, errors.New("User not found")
	}

	if user.Status == false {
		return nil, errors.New("User is disabled")
	}

	matched := PasswordMatch(user.Password, req.Password)
	if !matched {
		return nil, errors.New("Password didn't match")
	}

	token, errToken := CreateToken(*user, db)
	if errToken != nil {
		return nil, errors.New("Cannot create login token")
	}

	body := map[string]interface{}{
		"token": token,
	}
	return body, nil
}
