package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khanakia/jgo/pkg/core"
	"github.com/khanakia/jgo/pkg/dbc"
	"github.com/khanakia/jgo/pkg/logger"
	"github.com/khanakia/jgo/pkg/server"
	"github.com/khanakia/jgo/pkg/util"
	"github.com/pkg/errors"
)

const (
	RoleSaID     = 1
	RoleMemberID = 2
)

// User ...
type User struct {
	core.Model
	// Just a unique id for user so we will refer the user by UID not the ID as it's hard to guess and difficult to brute force
	UID       string `json:"-" gorm:"type:varchar(36)"`
	FullName  string `json:"fullName" gorm:"type:varchar(50)"`
	FirstName string `json:"firstName" gorm:"type:varchar(50)"`
	LastName  string `json:"lastName" gorm:"type:varchar(50)"`
	Email     string `json:"email" gorm:"type:varchar(100);unique"`
	Password  string `json:"-" gorm:"type:varchar(250)"`
	Secret    string `json:"-" gorm:"type:varchar(50)"` // Will be used for Login or Other function this will be internal and must never shared to frotend
	Status    bool   `json:"status" gorm:"type:boolean;default:true"`
	// ParentID         uint      `json:"parentId"`
	RoleID           uint      `json:"roleId"`
	Optin            bool      `json:"-" gorm:"type:boolean;default:false"`
	OptinAt          time.Time `json:"optinAt"`
	BusinessName     string    `json:"businessName" gorm:"type:varchar(50)"`
	Phone            string    `json:"phone" gorm:"type:varchar(25)"`
	CountryID        uint      `json:"countryId"`
	StateID          uint      `json:"stateId"`
	City             string    `json:"city" gorm:"type:varchar(50)"`
	PinCode          string    `json:"pinCode" gorm:"type:varchar(10)"`
	AddressLine1     string    `json:"addressLine1" gorm:"type:varchar(150)"`
	AddressLine2     string    `json:"addressLine2" gorm:"type:varchar(150)"`
	SignupSource     string    `json:"signupSource" gorm:"type:varchar(50)"`
	RandomPass       string    `json:"randomPass" gorm:"type:varchar(50)"`
	WelcomeEmailSent bool      `json:"welcomeEmailSent" gorm:"type:boolean;default:false"`

	// User Store Default Currency
	CurrencyCode string `json:"currencyCode" gorm:"type:varchar(3)"`
}

func (user *User) FillDefaults() {

	if user.RoleID == 0 {
		user.RoleID = RoleMemberID
	}

	if !user.Status {
		user.Status = true
	}

	if len(user.Password) <= 0 {
		randomPass := RandomPass()
		user.RandomPass = randomPass
		user.Password = GeneratePassword(randomPass)
	} else {
		user.Password = GeneratePassword(user.Password)
	}

	if len(user.Secret) <= 0 {
		// secret := uuid.New().String()
		// user.Secret = string(secret)
		user.Secret = util.GenerateUID()
	}

	if len(user.UID) <= 0 {
		// secret := uuid.New().String()
		// user.UID = string(secret)
		user.UID = util.GenerateUID()
	}

}

// func (user User) GetName() string {
// 	name := user.FirstName
// 	if len(user.LastName) > 0 {
// 		name = name + " " + user.LastName
// 	}
// 	return name
// }

// func GetAddress(user User, db *gorm.DB) uxmother.AddressView {
// 	return uxmother.GetFullAddress(user.AddressLine1, user.AddressLine2, user.CountryID, user.StateID, user.City, user.PinCode, db)
// }

// func (a User) ValidateBeforeInsert() ([]interface{}, error) {
// 	errorValidate := validation.ValidateStruct(&a,
// 		validation.Field(&a.Email, validation.Required, is.Email),
// 		// validation.Field(&a.ParentID, validation.Required),
// 		validation.Field(&a.RoleID, validation.Required),
// 		validation.Field(&a.Password, validation.Required),
// 	)
// 	data := util.ErroObjToArray(errorValidate)
// 	return data, errorValidate
// }

type Auth struct {
	Config
}

func (auth Auth) Version() string {
	return "0.01"
}

func (auth Auth) AutoMigrate() {
	auth.Dbc.Db.AutoMigrate(&User{})
}

type Config struct {
	Dbc    dbc.Dbc
	Logger logger.Logger
	// Hello        hello.Hello
	Server       server.Server
	EnableRoutes bool     `wire:"-"`
	Names        []string `wire:"-"`
}

func (config *Config) parse() {
	// set default values
	config.EnableRoutes = true

	// override values from config file
	err := util.ParseConfig("auth", &config)
	if err != nil {
		config.Logger.SugarLogger.Error(err)
		panic(err)
	}
	// fmt.Printf("%#v\n", config)
}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong11",
	})
}

func parseConfig() {
	student := Config{}

	err := util.ParseConfig("auth", &student)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v\n", student)
}

func New(config Config) Auth {
	config.parse()
	// fmt.Printf("%#v\n", config)

	// fmt.Println(hello.GetFname())

	// parseConfig()

	config.Server.Router.GET("/auth/ping", pingHandler)
	// p, err := config.Server.GetRouterGroup("private")
	// fmt.Println(err)
	// if err == nil {
	// 	fmt.Println(p)
	// }
	// p.GET("/a", pingHandler)

	auth := Auth{
		Config: config,
	}

	// auth.AutoMigrate()

	return auth
}

func (auth Auth) Register(user *User) (*User, error) {
	var err error
	db := auth.Dbc.Db

	user.FillDefaults()

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
