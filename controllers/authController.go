package controllers

import (
	"fmt"
	"net/http"
	"payment/database"
	"payment/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var data map[string]string

	// リクエストデータをパースする
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data["password"] != data["password_confirm"] {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Passwords do not match!"})
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: []byte(password),
	}

	fmt.Println("Saving user:", user)
	err = models.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal  error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
}

func Login(c *gin.Context) {
	var data map[string]string

	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	db, err := database.Connect()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	db.Where("email = ?", data["email"]).First(&user)
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"]))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// JWT
	claims := jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	// Cookie
	c.SetCookie("jwt", token, int(24*time.Hour.Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"jwt": token,
	})
}

type Claims struct {
	jwt.StandardClaims
}

func User(c *gin.Context) {
    // CookieからJWTを取得
    cookie, err := c.Cookie("jwt") // Loginで保存したもの
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "message": "unauthenticated",
        })
        return
    }

    // token取得
    token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte("secret"), nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{
            "message": "unauthenticated",
        })
        return
    }

    claims := token.Claims.(*jwt.StandardClaims)
    // User IDを取得
    id := claims.Issuer

    var user models.User
    db, err := database.Connect()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    db.Where("id = ?", id).First(&user)

    c.JSON(http.StatusOK, user)
}

func Logout(c *gin.Context) {
    cookie := http.Cookie{
        Name:     "jwt",
        Value:    "",                         // tokenを空にする
        Expires:  time.Now().Add(-time.Hour), // マイナス値を入れて期限切れ
        HttpOnly: true,
    }

    c.SetCookie(cookie.Name, cookie.Value, int(-time.Hour.Seconds()), cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
    c.Redirect(http.StatusTemporaryRedirect, "/login")
}
