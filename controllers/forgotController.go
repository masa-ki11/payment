package controllers

import (
	// "os"
	"math/rand"
	"net/http"
	// "net/smtp"
	"log"
	// "time"
	"fmt"
	"payment/database"
	"payment/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
)

func Forgot(c *gin.Context) {
	var data map[string]string

	// リクエストデータをパースする
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var user models.User
	if err := db.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		log.Printf("ユーザーの取得に失敗しました: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりませんでしたユーザー登録をしてください"})
		c.Abort()
		return
	}
	if user.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "管理者はパスワード変更できません。管理者にお問い合わせください。"})
		c.Abort()
		return
	}

	token := RandStringRunes(12)
	passwordReset := models.PasswordReset{
		Email: data["email"],
		Token: token,
	}

	// DBに保存
	db.Create(&passwordReset)

	 // メール送信用
	// // SMTPメール送信
	// from := os.Getenv("MAIL_FROM")
	// password := os.Getenv("MAIL_PASSWORD")
	// to := []string{
	// 	data["email"],
	// }
	// sendFrom := fmt.Sprintf("From: %s\n", from)
	// subject := fmt.Sprintf("Subject: %s\n", "Password Reset")
	// mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// url := "https://go-wallet-app.azurewebsites.net/reset-password?token=" + token
	// message := fmt.Sprintf("Click <a href=\"%s\">here</a> to reset password!", url)
	// auth := smtp.PlainAuth("", from, password, "smtp.alpha-prm.jp")
	// sendErr := smtp.SendMail(
	// 	"smtp.alpha-prm.jp:587",
	// 	auth,
	// 	from,
	// 	to,
	// 	[]byte(sendFrom+subject+mime+message),
	// )

	// if sendErr != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": sendErr.Error(),
	// 	})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "SUCCESS",
	// })


	// メール送信不可の代替案
	url := "/reset-password?token=" + token
	c.JSON(http.StatusOK, gin.H{
		"resetUrl": url,
	})
}

// ランダム文字列を返す
func RandStringRunes(n int) string {
	var lettersRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = lettersRunes[rand.Intn(len(lettersRunes))]
	}
	return string(b)
}

func Reset(c *gin.Context) {
	var data map[string]string

	// リクエストデータをパース
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// チェック
	if data["password"] != data["password_confirm"] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "パスワードが違います。",
		})
		return
	}

	var passwordReset = models.PasswordReset{}
	// JWT Tokenからデータを取得
	db, err := database.Connect()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
	err = db.Where("token = ?", data["token"]).Last(&passwordReset).Error
	fmt.Println(data["token"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid token!",
		})
		return
	}

	// パスワードをエンコード
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	db.Model(&models.User{}).Where("email = ?", passwordReset.Email).Update("password", password)

	c.JSON(http.StatusOK, gin.H{
		"message": "SUCCESS",
	})
}