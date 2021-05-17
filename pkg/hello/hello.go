package hello

type MyString string

var fname string

func SetFname(val string) {
	fname = val
}

func GetFname() string {
	return fname
}

type Hello struct {
	Name *MyString
}

func New() Hello {
	var name MyString
	name = "Aman"
	return Hello{
		Name: &name,
	}
}

func Sum(x int, y int) int {
	return x + y
}
