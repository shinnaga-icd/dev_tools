package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

// Menu はindexページのメニュー構造体型
type Menu struct {
	Name string
	URL  string
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	serve := gin.Default()
	serve.Use(gin.Logger())
	serve.LoadHTMLGlob("templates/*.tmpl.html")
	serve.Static("/static", "./static")

	// index page
	serve.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"menus": []Menu{{"リクエスト確認さん", "/request"}, {"暗号さん・複合さん", "/encrypt"}},
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

	serve.Run(":" + port)
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
