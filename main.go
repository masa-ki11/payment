package main

import (
	// "fmt"
	// "payment/database"
	"payment/routes"
	"time"
	// "net/http"
	"github.com/gin-gonic/gin"
	// "github.com/golang-migrate/migrate"
)

func main() {

		// タイムゾーンを JST (日本標準時) に設定する
		location, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			panic(err)
		}
		time.Local = location
	// サーバーをセットアップする
	r := gin.Default()

	// テンプレートをロードする
	r.LoadHTMLGlob("templates/*")
	r.Static("/js", "./js")
	r.Static("/static", "./static")
	// ルートを設定する
	routes.Setup(r)

	r.Run(":8080")

}
