package routes

import (
	"fmt"
	"net/http"
	"payment/controllers"
    "payment/middleware"
	"github.com/golang-jwt/jwt"
	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.GET("/logout", controllers.Logout)
	router.POST("/forgot", controllers.Forgot)
	router.POST("/reset", controllers.Reset)
    router.POST("/add-points", controllers.AddPoints)
    router.POST("/use-points", controllers.UsePoints)
    router.GET("/user-history", controllers.GetHistory)
    router.GET("/get-points", controllers.GetPoint)

    router.GET("/user", middleware.AdminMiddleware(), func(c *gin.Context) {
        // 現在のユーザー情報を取得
        user, err := controllers.GetCurrentUser(c)
        if err != nil {
            // エラー処理を行うか、デフォルト値を設定します
        }
        // 現在のユーザー情報をコンテキストに保存
        c.Set("currentUser", user)
        controllers.GetUser(c)
    })
    router.GET("/all-users", middleware.AdminMiddleware(), func(c *gin.Context) {
        // 現在のユーザー情報を取得
        user, err := controllers.GetCurrentUser(c)
        if err != nil {
            // エラー処理を行うか、デフォルト値を設定します
        }
        // 現在のユーザー情報をコンテキストに保存
        c.Set("currentUser", user)

        controllers.GetAllUsers(c)
    })
	router.GET("/", AuthRequired, func(c *gin.Context) {
        user, err := controllers.GetCurrentUser(c)
        if err != nil {
            // エラー処理を行うか、デフォルト値を設定します
        }
        c.HTML(200, "home.html", gin.H{
            "admin": user.Admin,
            "user": user,
        })
    })
    router.GET("/history", AuthRequired, func(c *gin.Context) {
        user, err := controllers.GetCurrentUser(c)
        if err != nil {
            // エラー処理を行うか、デフォルト値を設定します
        }
        c.HTML(200, "history.html", gin.H{
            "admin": user.Admin,
            "user": user,
        })
    })
	router.GET("/register", func(c *gin.Context) {
        c.HTML(200, "register.html", nil)
    })
    router.GET("/login", func(c *gin.Context) {
        c.HTML(200, "login.html", nil)
    })
	router.GET("/forgot-password", func(c *gin.Context) {
        c.HTML(200, "forgot-password.html", nil)
    })
	router.GET("/reset-password", func(c *gin.Context) {
        c.HTML(200, "reset-password.html", nil)
    })

}

func AuthRequired(c *gin.Context) {
    cookie, err := c.Cookie("jwt")
	fmt.Println(cookie)
    if err != nil {
        c.Redirect(http.StatusTemporaryRedirect, "/login")
        c.Abort()
        return
    }

    token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte("secret"), nil
    })

    if err != nil || !token.Valid {
        c.Redirect(http.StatusTemporaryRedirect, "/login")
        c.Abort()
        return
    }
}
