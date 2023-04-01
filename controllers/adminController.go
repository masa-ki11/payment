
package controllers

import (
	"fmt"
	// "net/http"
	"strconv"
	"payment/database"
	"payment/models"
	"github.com/golang-jwt/jwt"
	"github.com/gin-gonic/gin"
)

// func HomeWithAdmin(c *gin.Context) {
// 	jwtToken, err := c.Cookie("jwt") // Loginで保存したもの
//     if err != nil {
//         c.JSON(http.StatusUnauthorized, gin.H{
//             "message": "unauthenticated",
//         })
//         return
//     }
// 	token, err := jwt.ParseWithClaims(jwtToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
//         return []byte("secret"), nil
//     })
// 	var currentUserID int
// 	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
//         // ユーザーIDを取得
//         currentUserID, err = strconv.Atoi(claims.Issuer)
//         if err != nil {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
//             return
//         }
//     } else {
//         c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
//         return
//     }
// 	db, err := database.Connect()
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
//         return
//     }

//     var user models.User
//     db.Where("id = ?", currentUserID).First(&user)

//     // ユーザーが管理者であるかどうかを確認
//     if user.Admin {
//         c.HTML(http.StatusOK, "all-users.html", nil)
//     } else {
//         c.Redirect(http.StatusFound, "/") // home.htmlにリダイレクト
//     }
// }

func GetCurrentUser(c *gin.Context) (*models.User, error) {
	jwtToken, err := c.Cookie("jwt")
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(jwtToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		userID, err := strconv.Atoi(claims.Issuer)
		if err != nil {
			return nil, err
		}

		db, err := database.Connect()
		if err != nil {
			return nil, err
		}

		var user models.User
		db.Where("id = ?", userID).First(&user)

		return &user, nil
	}

	return nil, fmt.Errorf("invalid token")
}
