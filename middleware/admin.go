package middleware

import (
	// "fmt"
	"net/http"
	"strconv"
	"payment/database"
	"payment/models"
	// "payment/controllers"
	"github.com/golang-jwt/jwt"
	"github.com/gin-gonic/gin"
)


func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        jwtToken, err := c.Cookie("jwt") // Loginで保存したもの
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "message": "unauthenticated",
            })
            c.Abort() // ミドルウェアチェーンを中断
            return
        }
        token, err := jwt.ParseWithClaims(jwtToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte("secret"), nil
        })
        var currentUserID int
        if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
            // ユーザーIDを取得
            currentUserID, err = strconv.Atoi(claims.Issuer)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
                c.Abort()
                return
            }
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
            c.Abort()
            return
        }
        db, err := database.Connect()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
            c.Abort()
            return
        }

        var user models.User
        db.Where("id = ?", currentUserID).First(&user)

        // ユーザーが管理者であるかどうかを確認
        if !user.Admin {
            c.Redirect(http.StatusFound, "/") // home.htmlにリダイレクト
            c.Abort()
            return
        }

    }
}
