package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/errgroup"
	"html/template"
	"net/http"
	"time"
)

var db = make(map[string]string)

var html = template.Must(template.New("https").Parse(`
<html>
<head>
	<title>Https Test</title>
	<script src="/assets/app.js"></script>
</head>
<body>
	<h1 style="color:red;">Welcome, Ginner!</h1>
</body>
</html>
`))

var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@quapa.com", "phone": "234"},
}

type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type Login struct {
	User     string `form:"user" json:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type Person struct {
	Name    string `form:"name"`
	Address string `form:"address"`
}

type formA struct {
	Foo string `json:"foo" xml:"foo" binding:"required"`
}

type formB struct {
	Bar string `json:"bar" xml:"bar" binding:"required"`
}

type myForm struct {
	Colors []string `form:"colors[]"`
}

type PersonWithUri struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

type structA struct {
	FieldA string `form:"field_a"`
}

type StructB struct {
	NestedStruct structA
	FieldB       string `form:"field_b"`
}

type StructC struct {
	NestedStructPointer *structA
	FieldC              string   `form:"field_c"`
}

type StructD struct {
	NestedAnonyStruct struct {
		FieldX string `form:"field_x"`
	}
	FieldD string `form:"field_d"`
}

// Booking 包含绑定和验证的数据
type Booking struct {
	CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn,bookabledate" time_format:"2006-01-02"`
}

var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

var g errgroup.Group

func router01() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 01",
			},
		)
	})

	return e
}

func router02() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 02",
			},
		)
	})

	return e
}
