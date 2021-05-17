package auth

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (a User) ValidateRegister() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.FirstName, validation.Required),
		// validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.LastName, validation.Required),
	)
}

// Register - This is the most basic function to register user all you need to pass is EMAIL and ParentID
// All other params are fully optional
func Register(user *User, db *gorm.DB) (*User, error) {
	var err error
	// YTD
	// user.FillDefaults()

	// YTD
	// _, err := user.ValidateBeforeInsert()
	// if err != nil {
	// 	return nil, err
	// }

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	isEmailExists := CheckEmailExists(user.Email, db)
	if isEmailExists {
		return nil, errors.New("Email already exists")
	}

	// parent := GetUser(user.ParentID, db)
	// if parent == nil {
	// 	return nil, errors.New("Parent requried")
	// }

	err = db.Create(user).Error
	if err != nil {
		return nil, errors.New("Server error. Please contact support")
	}

	return user, nil
}
