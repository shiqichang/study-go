package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"github.com/go-playground/validator/v10"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

// 自定义中间件

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Set("example", "12345")

		// 请求前
		c.Next()

		// 请求后
		latency := time.Since(t)
		log.Print(latency)

		// 获取发送的 status
		status := c.Writer.Status()
		log.Println(status)
	}
}

func setupRouter() *gin.Engine {
	//r := gin.Default()
	r := gin.New()
	//r.Use(Logger())

	// 自定义日志文件
	r.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				params.ClientIP,
				params.TimeStamp.Format(time.RFC1123),
				params.Method,
				params.Path,
				params.Request.Proto,
				params.StatusCode,
				params.Latency,
				params.Request.UserAgent(),
				params.ErrorMessage,
		)
	}))

	// 记录日志
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色
	//gin.DisableConsoleColor()
	// 强制日志颜色化
	gin.ForceConsoleColor()

	// 记录到文件
	f, _ := os.Create("./log/gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)
	// 同时将日志写入文件和控制台
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 定义路由日志的格式
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	r.GET("/test", func(c *gin.Context) {
		example := c.MustGet("example").(string)
		log.Println(example)
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	//authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
	//	"foo":  "bar", // user:foo password:bar
	//	"manu": "123", // user:manu password:123
	//}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	//authorized.POST("admin", func(c *gin.Context) {
	//	user := c.MustGet(gin.AuthUserKey).(string)
	//
	//	// Parse JSON
	//	var json struct {
	//		Value string `json:"value" binding:"required"`
	//	}
	//
	//	if c.Bind(&json) == nil {
	//		db[user] = json.Value
	//		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	//	}
	//})

	// AsciiJSON
	r.GET("/some-json", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}
		c.AsciiJSON(http.StatusOK, data)
	})

	// HTML 渲染
	r.LoadHTMLGlob("templates/**/*")
	r.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})
	r.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
		})
	})

	// 自定义模板功能
	r.Delims("{[{", "}]}")
	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	r.LoadHTMLFiles("templates/raw.tmpl")
	r.GET("/raw", func(c *gin.Context) {
		c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{}{
			"now": time.Date(2022, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	// HTTP2 server 推送
	r.Static("/assets", "./assets")
	r.SetHTMLTemplate(html)

	r.GET("/", func(c *gin.Context) {
		time.Sleep(5*time.Second)
		if pusher := c.Writer.Pusher(); pusher != nil {
			// 使用 pusher.Push() 做服务器推送
			if err := pusher.Push("/assets/app.js", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
		c.HTML(http.StatusOK, "https", gin.H{
			"status": "success",
		})
	})

	// JSONP
	r.GET("/jsonp", func(c *gin.Context) {
		data := map[string]interface{}{
			"foo": "bar",
		}

		// /JSONP?callback=x
		c.JSONP(http.StatusOK, data)
	})

	// Multipart/Urlencoded 绑定
	r.POST("/login", func(c *gin.Context) {
		var form LoginForm
		if c.ShouldBind(&form) == nil {
			if form.User == "user" && form.Password == "password" {
				c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(http.StatusOK, gin.H{"status": "unauthorized"})
			}
		}
	})

	// Multipart/Urlencoded 表单
	r.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(http.StatusOK, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})

	// PureJSON
	// 提供 unicode 实体
	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"html": "<b>hello, world!</b>",
		})
	})
	// 提供字面字符
	r.GET("/purejson", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, gin.H{
			"html": "<b>hello, world!</b>",
		})
	})

	// Query 和 post form
	r.POST("/post", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")

		fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})

	// SecureJSON
	r.GET("/secure-json", func(c *gin.Context) {
		names := []string{"lena", "austin", "foo"}
		c.SecureJSON(http.StatusOK, names)
	})

	// XML/JSON/YAML/ProtoBuf 渲染
	r.GET("/some-json/v1", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hey",
			"status":   http.StatusOK,
		})
	})
	r.GET("/some-json/v2", func(c *gin.Context) {
		var msg struct{
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		c.JSON(http.StatusOK, msg)
	})

	r.GET("/some-xml", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{
			"message": "hey",
			"status":  http.StatusOK,
		})
	})
	r.GET("/some-yaml", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{
			"message": "hey",
			"status":  http.StatusOK,
		})
	})
	r.GET("/some-protobuf", func(c *gin.Context) {
		reps := []int64{1, 2}
		label := "test"
		data := &protoexample.Test{
			Label: &label,
			Reps:  reps,
		}
		c.ProtoBuf(http.StatusOK, data)
	})

	// 上传文件
	// 单文件
	r.MaxMultipartMemory = 8 << 20
	r.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		dst := "./" + file.Filename
		c.SaveUploadedFile(file, dst)
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	// 多文件
	r.POST("/multi-upload", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)
			dst := "./" + file.Filename
			c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	// 不使用默认的中间件
	//r := gin.New

	// 从 reader 读取数据
	r.GET("/some-data-from-reader", func(c *gin.Context) {
		response, err := http.Get("http://img-ys011.didistatic.com/static/gungnir_data_ingestion_material_img/3d9fbe3e-6ab0-46b1-a3ba-21a2fd90b245")
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")
		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}
		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})

	// 使用 BasicAuth 中间件
	// 路由组使用 gin.BasicAuth() 中间件
	// gin.Accounts 是 map[string]string 的一种快捷方式
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello",
		"manu":   "4321",
	}))

	authorized.GET("/secrets", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	// 使用 HTTP 方法
	//r.GET("/get", get)
	//r.POST("/post", post)
	//r.PUT("/put", put)
	//r.DELETE("/delete", delete)
	//r.PATCH("/patch", patch)
	//r.HEAD("/head", head)
	//r.OPTIONS("/options", options)

	// 使用中间件
	//r := gin.New()
	//r.Use(gin.Logger())
	//r.Use(gin.Recovery())
	//r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// 只绑定 url 查询字符串
	r.Any("/testing", startPage)

	// 在中间件中使用 Goroutine
	r.GET("/long_async", func(c *gin.Context) {
		// 创建在 goroutine 中使用的副本
		cCp := c.Copy()
		go func() {
			time.Sleep(5 * time.Second)
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})

	r.GET("/long_sync", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		log.Println("Done! in path " + c.Request.URL.Path)
	})

	// 支持 Let's Encrypt
	//m := autocert.Manager{
	//	Prompt:     autocert.AcceptTOS,
	//	HostPolicy: autocert.HostWhitelist("example1.com", "example2.com"),
	//	Cache:      autocert.DirCache("/var/www/.cache"),
	//}
	//log.Fatal(autotls.RunWithManager(r, &m))

	// 映射查询字符串或表单参数
	r.POST("/post/v1", func(c *gin.Context) {
		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")
		fmt.Printf("ids: %v; names: %v", ids, names)
	})

	// 查询字符串参数
	r.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")
		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	// 模型绑定和验证
	// 绑定 JSON
	r.POST("/login-json", func(c *gin.Context) {
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.User != "manu" || json.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})
	// 绑定 XML
	// <?xml version="1.0" encoding="UTF-8"?>
	//	//	<root>
	//	//		<user>manu</user>
	//	//		<password>123</password>
	//	//	</root>
	r.POST("/login-xml", func(c *gin.Context) {
		var xml Login
		if err := c.ShouldBindXML(&xml); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if xml.User != "manu" || xml.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})
	// 绑定 HTML 表单
	r.POST("/login-form", func(c *gin.Context) {
		var form Login
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if form.User != "manu" || form.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// 绑定 HTML 复选框
	r.POST("/bind-html", formHandler)

	// 绑定 Uri
	r.GET("/:name/:id", func(c *gin.Context) {
		var person PersonWithUri
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"name": person.Name, "uuid": person.ID})
	})

	// 绑定表单数据至自定义结构体
	// 仅支持没有 form 的嵌套结构体
	r.GET("/getb", GetDataB)
	r.GET("/getc", GetDataC)
	r.GET("/getd", GetDataD)

	// 自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
	}
	r.GET("/bookable", getBookable)

	// 设置和获取 cookie
	r.GET("/cookie", func(c *gin.Context) {
		cookie, err := c.Cookie("gin_cookie")
		if err != nil {
			cookie = "NoSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}
		fmt.Printf("Cookie value: %s \n", cookie)
	})

	// 路由参数
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	// 路由组
	//v1 := r.Group("/v1")
	//{
	//	v1.POST("/login", login)
	//	v1.POST("/submit", submit)
	//	v1.POST("/read", read)
	//}

	// 重定向
	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
	})
	r.POST("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ping")
	})
	r.GET("/test/v1", func(c *gin.Context) {
		c.Request.URL.Path = "/test/v2"
		r.HandleContext(c)
	})
	r.GET("/test/v2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})

	// 静态文件服务
	//r.Static("/assets", "./assets")
	//r.StaticFS("/more_static", http.Dir("my_first_system"))
	//r.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// 静态文件嵌入
	//t, err := loadTemplate()
	//if err != nil {
	//	panic(err)
	//}
	//r.SetHTMLTemplate(t)
	//r.GET("/static", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "/html/index.tmpl", nil)
	//})

	return r
}

func startPage(c *gin.Context) {
	var person Person
	if c.ShouldBindQuery(&person) == nil {
		log.Println("===== Only Bind By Query String =====")
		log.Println(person.Name)
		log.Println(person.Address)
	}
	c.String(http.StatusOK, "Success")
}

func SomeHandler(c *gin.Context) {
	objA := formA{}
	objB := formB{}

	// c.ShouldBind 使用 c.Request.Body，不可重用
	if errA := c.ShouldBind(&objA); errA != nil {
		c.String(http.StatusOK, `the body should be formA`)
	} else if errB := c.ShouldBind(&objB); errB != nil {
		c.String(http.StatusOK, `the body should be formB`)
	} else {

	}

	// 多次绑定，只有某些格式需要，如 JSON, XML, MsgPack, ProtoBuf. 对于其他格式，如 Query, Form, FormPost, FormMultipart 可以多次调用 c.ShouldBind()
	if errA := c.ShouldBindBodyWith(&objA, binding.JSON); errA == nil {
		c.String(http.StatusOK, `the body should be formA`)
	} else if errB := c.ShouldBindBodyWith(&objB, binding.JSON); errB == nil {
		c.String(http.StatusOK, `the body should be formB JSON`)
	} else if errB2 := c.ShouldBindBodyWith(&objB, binding.XML); errB2 == nil {
		c.String(http.StatusOK, `the body should be formB XML`)
	} else {

	}
}

func formHandler(c *gin.Context) {
	var fakeForm myForm
	c.ShouldBind(&fakeForm)
	c.JSON(http.StatusOK, gin.H{"color": fakeForm.Colors})
}

func GetDataB(c *gin.Context) {
	var b StructB
	c.Bind(&b)
	c.JSON(http.StatusOK, gin.H{"a": b.NestedStruct, "b": b.FieldB})
}

func GetDataC(c *gin.Context) {
	var b StructC
	c.Bind(&b)
	c.JSON(http.StatusOK, gin.H{"a": b.NestedStructPointer, "c": b.FieldC})
}

func GetDataD(c *gin.Context) {
	var b StructD
	c.Bind(&b)
	c.JSON(http.StatusOK, gin.H{"x": b.NestedAnonyStruct, "d": b.FieldD})
}

func getBookable(c *gin.Context) {
	var b Booking
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// loadTemplate 加载由 go-assets-builder 嵌入的模板
//func loadTemplate() (*template.Template, error) {
//	t := template.New("")
//	for name, file := range Assets.Files {
//		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
//			continue
//		}
//		h, err := ioutil.ReadAll(file)
//		if err != nil {
//			return nil, err
//		}
//		t, err = t.New(name).Parse(string(h))
//		if err != nil {
//			return nil, err
//		}
//	}
//	return t, nil
//}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	//r.Run(":8080")
	//r.Run(":8080", "./testdata/server.pem", "./testdata/server.key")

	// 优雅地重启或停止
	// 自定义 HTTP 配置
	srv := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		// 服务链接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

	// 运行多个服务
	//server01 := &http.Server{
	//	Addr: ":8080",
	//	Handler: router01(),
	//	ReadTimeout: 5 * time.Second,
	//	WriteTimeout: 10 * time.Second,
	//}
	//server02 := &http.Server{
	//	Addr: ":8081",
	//	Handler: router02(),
	//	ReadTimeout: 5 * time.Second,
	//	WriteTimeout: 10 * time.Second,
	//}
	//
	//g.Go(func() error {
	//	return server01.ListenAndServe()
	//})
	//g.Go(func() error {
	//	return server02.ListenAndServe()
	//})
	//
	//if err := g.Wait(); err != nil {
	//	log.Fatal(err)
	//}
}
