package mail

import (
	"bytes"
	"html/template"
	"net/http"
	"net/smtp"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const Version = "0.0.1"

type Mail struct{}

type Message struct {
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
}

// New initialize
func New() *Mail {
	mail := &Mail{}
	return mail
}

// https://stackoverflow.com/questions/44675087/golang-template-variable-isset
// {{if (avail "Name" .Data)}}Name is: {{.Data.Name}}{{else}}Name is unavailable.{{end}}
func avail(name string, data interface{}) bool {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}
	return v.FieldByName(name).IsValid()
}

// {{ hello }}
func hello() string {
	return "Hello"
}

// Compile - Template files
func Compile(files []string, data interface{}) (string, error) {

	buf := new(bytes.Buffer)

	var paths []string
	for _, v := range files {
		// fmt.Printf("2**%d = %d\n", i, v)
		paths = append(paths, "./resources/emails/"+v)
	}

	// files := []string{
	// 	"./resources/views/register.html",
	// 	"./resources/views/base1.html",
	// }

	// t, errP := (template.ParseFiles(paths...))
	t := template.Must(template.New(files[0]).Funcs(template.FuncMap{
		"htmlSafe": func(html string) template.HTML {
			return template.HTML(html)
		},
		"hello": hello,
		"avail": avail,
	}).ParseFiles(paths...))

	// if errP != nil {
	//
	// 	return "", errP
	// }

	// data := struct {
	// 	AppName string
	// 	Data    interface{}
	// }{
	// 	AppName: viper.GetString("app_name"),
	// 	Data:    templateData,
	// }
	err := t.Execute(buf, data)
	if err != nil {

		return "", err
	}

	body := buf.String()

	return body, nil

}

func Send(msg *Message) error {
	// testSend()
	from := viper.GetString("mail.from")
	host := viper.GetString("mail.host")
	port := viper.GetString("mail.port")
	user := viper.GetString("mail.user")
	pass := viper.GetString("mail.password")

	addr := host + ":" + port
	auth := smtp.PlainAuth("", user, pass, host)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// Merge the To, Cc, and Bcc fields
	to := make([]string, 0, len(msg.To)+len(msg.Cc)+len(msg.Bcc))
	to = append(append(append(to, msg.To...), msg.Cc...), msg.Bcc...)

	body := []byte(
		"To: " + strings.Join(msg.To, ",") + "\r\n" +
			"Cc: " + strings.Join(msg.Cc, ",") + "\r\n" +
			"Bcc: " + strings.Join(msg.Bcc, ",") + "\r\n" +
			"Subject: " + msg.Subject + "\r\n" +
			mime +
			"\r\n" +
			msg.Body +
			"\r\n")
	err := smtp.SendMail(addr, auth, from, to, body)
	if err != nil {
		// log.Fatal(err)
		return err
	}

	return nil
}

func TestSendHandler(c *gin.Context) {
	// c.JSON(200, gin.H{
	// 	"message": "pong",
	// })
	files := []string{
		"register.html",
		"base.html",
	}

	templateData := struct {
		Name string
		URL  string
	}{
		Name: "Dhanush",
		URL:  "http://geektrust.in",
	}
	body, _ := Compile(files, templateData)

	msg := &Message{
		To:      []string{"khanakia@gmail.com", "khanakia1@gmail.com"},
		Cc:      []string{"cc1@khanakia.com"},
		Bcc:     []string{"bcc1@khanakia.com"},
		Subject: "Testing",
		Body:    body,
	}
	Send(msg)

	// c.HTML(200, body, "")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(body))
}
