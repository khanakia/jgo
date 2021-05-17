package auth

import (
	"math/rand"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/khanakia/jgo/pkg/util"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GenerateUID() string {
	secret := uuid.New().String()
	return secret
}

func GetUserSecret(user User, db *gorm.DB) string {
	if len(user.Secret) == 0 {
		user.Secret = util.GenerateUID()

		db.Save(&user)
		return user.Secret
	}
	return user.Secret
}

func GetUserUID(user User, db *gorm.DB) string {
	if user.UID == "" {
		user.UID = util.GenerateUID()
		db.Save(&user)
		return user.UID
	}
	return user.UID
}

/*
 * This signature is used to create a JWT token
 * Benefits - By adding the user secret with the appsecret it makes the app more secure
 * Let say if appSecret is compromised then stil nobody can generate tokens without user secret
 * If userSecret compromised it will affect only single user not all the users
 * If user wants to focefully logout for all the applications we simply update his userSecret
 * FUTURE CONSIDERATION - Add jwt to token to the blacklist if users logout
 */
func GetSignature(user User, db *gorm.DB) string {
	userSecret := GetUserSecret(user, db)

	appSecret := util.GetEnv("appSecret", "IBIrewORShiVReBASTer")
	signature := appSecret + ":" + userSecret
	return signature
}

// GeneratePassword - Create Bcrypt from string
func GeneratePassword(password string) string {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash)
}

// PasswordMatch - Compare two passwords are equal
func PasswordMatch(password string, password1 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(password1))
	if err == nil {
		return true
	}
	return false
}

func RandomPass() string {
	rand.Seed(time.Now().UnixNano())

	chars := []rune("0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // E.g. "ExcbsVQs"
	return str
}

func CreateToken(user User, db *gorm.DB) (string, error) {
	expirationTime := time.Now().Add(500 * time.Minute) // 500 minute

	claims := &Claims{
		// ID:    user.ID,
		Email: user.Email,
		UID:   user.UID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signature := []byte(GetSignature(user, db))
	tokenString, err := token.SignedString(signature)

	return tokenString, err
}

func GetUserFromContext(c *gin.Context) (User, error) {
	userc, _ := c.Get("user")

	if userc == nil {
		return User{}, errors.New("User not found")
	}
	user, _ := userc.(User)
	return user, nil
}

func CheckEmailExists(email string, db *gorm.DB) bool {
	var count int64
	db.Model(User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return true
	}

	return false
}

func FindByEmail(email string, db *gorm.DB) *User {
	var user User
	email = strings.ToLower(email)
	res := db.First(&user, &User{Email: email})

	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return &user
}

func GetUser(id uint, db *gorm.DB) *User {
	var user User
	res := db.First(&user, id)
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return &user
}
