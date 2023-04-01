package main

import (
	// "fmt"
	// "payment/database"
	"payment/routes"

	// "net/http"
	"github.com/gin-gonic/gin"
	// "github.com/golang-migrate/migrate"
)

func main() {

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
