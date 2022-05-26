package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/jung-kurt/gofpdf"
)

// Menu はindexページのメニュー構造体型
type Menu struct {
	Name string
	URL  string
}

// Controller
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	serve := gin.Default()
	serve.Use(gin.Logger())
	serve.LoadHTMLGlob("templates/*.tmpl.html")
	serve.Static("/static", "./static")

	store := cookie.NewStore([]byte("secret"))
	serve.Use(sessions.Sessions("mysession", store))

	// index page
	serve.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{"menus": []Menu{
			{"リクエスト確認さん", "/request"},
			{"暗号さん・複合さん", "/encrypt"},
			{"PDF出し太郎", "/pdf"},
			{"SESSION入り太郎", "/login"},
			{"SESSION見太郎", "/session"},
			{"ダイアログ確認したい君", "/dialog"}
		},
		})
	})

	// request page
	serve.GET("/request", func(c *gin.Context) {
		c.HTML(http.StatusOK, "request.tmpl.html", gin.H{
			"message": "GET Reqest",
		})
	})
	serve.POST("/request", func(c *gin.Context) {
		p1 := "p4"
		p2 := "bhreq"

		c.HTML(http.StatusOK, "request.tmpl.html", gin.H{
			"message": "POST Reqest",
			p1:        c.PostForm(p1),
			p2:        c.PostForm(p2),
		})
	})

	// encrypt page
	serve.GET("/encrypt", func(c *gin.Context) {
		c.HTML(http.StatusOK, "encrypt.tmpl.html", gin.H{})
	})
	serve.POST("/encrypt", func(c *gin.Context) {
		var result string
		var err error
		text := c.PostForm("text")
		key := c.PostForm("key")
		iv := c.PostForm("iv")

		if c.PostForm("encrypt") == "1" {
			result, err = encrypt(text, key, iv)
		} else {
			result, err = decrypt(text, key, iv)
		}

		c.HTML(http.StatusOK, "encrypt.tmpl.html", gin.H{
			"result": result,
			"text":   text,
			"key":    key,
			"iv":     iv,
			"error":  err,
		})
	})

	// outputpdf
	serve.GET("/pdf", func(c *gin.Context) {
		download := c.Query("download")
		if download == "pdf" {
			downloadPdf(c.Writer, c.Request)
		}

		c.HTML(http.StatusOK, "outputpdf.tmpl.html", gin.H{})
	})
	serve.POST("/pdf", func(c *gin.Context) {
		downloadPdf(c.Writer, c.Request)
	})

	//login session
	serve.POST("/login", func(c *gin.Context) {
		param := "bhreq"
		session := sessions.Default(c)

		session.Set(param, c.PostForm(param))
		session.Save()

		c.HTML(http.StatusOK, "login.tmpl.html", gin.H{
			"message": "session saved",
			param:     session.Get(param),
		})

	})

	//session confirm
	serve.GET("/session", func(c *gin.Context) {
		param := "bhreq"
		session := sessions.Default(c)
		sessionVal := session.Get(param)

		// session.Clear()
		// session.Save()

		c.HTML(http.StatusOK, "session.tmpl.html", gin.H{
			"message": "get session",
			param:     sessionVal,
		})

	})

	// dialog page
	serve.GET("/dialog", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dialog.tmpl.html", gin.H{})
	})
	serve.POST("/dialog", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dialog.tmpl.html", gin.H{})
	})

	serve.Run(":" + port)
}

// download PDF
func downloadPdf(w http.ResponseWriter, r *http.Request) {
	//make pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello world")
	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		log.Printf("err\n")
	}
	//download
	w.Header().Set("Content-Disposition", "attachment; filename=hello.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, "hello.pdf")
}

//  暗号化 AES
func encrypt(plainText string, key string, iv string) (string, error) {

	hexKey, _ := hex.DecodeString(key)
	hexIv, _ := hex.DecodeString(iv)

	block, err := aes.NewCipher([]byte(hexKey))
	if err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, []byte(hexIv))
	cipherText := make([]byte, len([]byte(plainText)))
	cfb.XORKeyStream(cipherText, []byte(plainText))

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// 復号化 AES
func decrypt(encrypted string, key string, iv string) (string, error) {
	byteEnc, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	cfbdec := cipher.NewCFBDecrypter(block, []byte(iv))
	cipherText := make([]byte, len(byteEnc))
	cfbdec.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

