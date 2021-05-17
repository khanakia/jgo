// Package jutil provide some common sets of shared function
package util

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type H map[string]interface{}

// Define n as number to limit the length of the random string
func GenerateRandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Generates random uuid string
func GenerateUID() string {
	secret := uuid.New().String()
	return secret
}

// Generates MDF of random uuid string
func GenerateMd5Uid() string {
	secret := uuid.New().String()
	key := []byte(secret)
	hash := md5.Sum(key)
	return hex.EncodeToString(hash[:])
}

// ModelType get value's model type
func ModelType(value interface{}) reflect.Type {
	reflectType := reflect.Indirect(reflect.ValueOf(value)).Type()

	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}

	return reflectType
}

// Use Make the unused value used so golang will not give error while compiling
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

// GetEnv get key environment variable if exist otherwise return defalutValue
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

type ErrorObject struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Idx     int    `json:"idx"`
}

// Convert github.com/go-ozzo/ozzo-validation errors to proper array of errors for the JSON
func ErroObjToArray(errors error) []interface{} {
	if errors == nil {
		return []interface{}{}
	}

	var errs []interface{}

	i := 0
	for field, v := range errors.(validation.Errors) {

		message := strings.Title(field) + " " + v.Error()
		errs = append(errs, &ErrorObject{
			Field:   field,
			Message: message,
			Idx:     i,
		})

		i++
	}

	return errs
}

/*
  https://gqlgen.com/reference/errors/
*/

func OzzoErrToGraphqlErrors(errors error, ctx context.Context) {
	i := 0
	for field, v := range errors.(validation.Errors) {
		message := strings.Title(field) + " " + v.Error()
		graphql.AddError(ctx, &gqlerror.Error{
			Path:    graphql.GetPath(ctx),
			Message: message,
			Extensions: map[string]interface{}{
				"idx":   i,
				"field": field,
			},
		})
		i++
	}
}

// Throw Malfunction error for Gin Request
func GinErrorMalfunction(c *gin.Context) {
	error := &ResponseError{
		Message: "Malfunctioned request.",
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": error})
}

// Throw Bad Request error for Gin Request
func GinErrorBadRequest(c *gin.Context) {
	error := &ResponseError{
		Message: "Bad request.",
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": error})
}

// Print any interface to type JSON
func PrintToJSON(val interface{}) {
	b, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

}

// https://gist.github.com/stoewer/fbe273b711e6a06315d19552dd4d33e6
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// Convert any string to snakecase StudentID will become student_id
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.Replace(snake, "__", "_", -1)

	return strings.ToLower(snake)
}

// Limit the decimal value e.g. 2.34343434 will become 2.34 if precision is set to 2
func FloatPrecision(num float64, precision int) float64 {
	p := math.Pow10(precision)
	value := float64(int(num*p)) / p
	return value
}

// https://golangcode.com/check-if-element-exists-in-slice/
// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func GetIndexFromStringSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func GetIndexFromUintSlice(slice []uint, val uint) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func GetFullName(firstName string, lastName string) string {
	name := firstName
	if len(lastName) > 0 {
		name = name + " " + lastName
	}
	return name
}

func ParseConfig(name string, config interface{}) error {
	configmap := viper.GetStringMap(name)
	jsonbody, err := json.Marshal(configmap)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonbody, &config); err != nil {
		return err
	}

	return nil
}
