package controllers

import (
	// "fmt"
	"errors"
	"net/http"
	"payment/database"
	"payment/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserWithPoints struct {
    models.User `json:"user"`
    Point uint `json:"point"`
}
func GetAllUsers(c *gin.Context) {
    db, err := database.Connect()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    var usersWithPoints []UserWithPoints
    db.Table("users").
        Select("users.*, points.point").
        Joins("left join points on users.id = points.user_id").
        Where("users.delete_flag = 0").
        Scan(&usersWithPoints)

    // コンテキストから現在のユーザー情報を取得
    currentUser, exists := c.Get("currentUser")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    user := currentUser.(*models.User)

    c.HTML(http.StatusOK, "all_users.html", gin.H{
        "users": usersWithPoints,
        "admin": user.Admin,
    })
}


func GetUser(c *gin.Context) {
    id := c.Query("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
        return
    }

    db, err := database.Connect()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    var user models.User
    db.Where("id = ?", id).First(&user)

    // コンテキストから現在のユーザー情報を取得
    currentUser, exists := c.Get("currentUser")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    current := currentUser.(*models.User)

    if user.ID == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.HTML(http.StatusOK, "user.html", gin.H{
        "user": user,
        "admin": current.Admin,
    })
}

func AddPoints(c *gin.Context) {
    // リクエストボディからユーザーIDとポイント数を取得する
    var req struct {
        UserIDs []string `json:"user_ids"`
        Points  uint     `json:"points"`
    }

    if err := c.BindJSON(&req); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }
    db, err := database.Connect()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    // ユーザーIDごとにポイントを加算する
    for _, userIDStr := range req.UserIDs {
        userID64, err := strconv.ParseUint(userIDStr, 10, 32)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user_id value"})
            return
        }
        userID := uint(userID64)
        var point models.Point
        if err := db.Where("user_id = ?", userID).First(&point).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                // ポイントレコードが存在しない場合は新規作成する
                point = models.Point{
                    UserID: userID,
                    Point:  req.Points,
                }
                if err := db.Create(&point).Error; err != nil {
                    c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create point record"})
                    return
                }
            } else {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve point record"})
                return
            }
        } else {
            // ポイントレコードが存在する場合は加算する
            point.Point += req.Points
            if err := db.Save(&point).Error; err != nil {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update point record"})
                return
            }
        }
        pointHistory := models.PointHistory{
            UserID:    userID,
            Point:     req.Points,
            Action:    "付与",
            Details:   "ポイント付与",
            CreatedAt: time.Now(),
        }
        if err := db.Create(&pointHistory).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create point history record"})
            return
        }
    }

    c.JSON(http.StatusOK, gin.H{})
}

func UsePoints(c *gin.Context) {
    // リクエストボディからユーザーIDとポイント数、使用用途を取得する
    var req struct {
        UserIDs []string `json:"user_ids"`
        Points  uint   `json:"points"`
        Details string `json:"details"`
    }

    if err := c.BindJSON(&req); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    db, err := database.Connect()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    for _, userIDStr := range req.UserIDs {
        userID64, err := strconv.ParseUint(userIDStr, 10, 32)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user_id value"})
            return
        }
        userID := uint(userID64)

        var point models.Point
        if err := db.Where("user_id = ?", userID).First(&point).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve point record"})
            return
        }

        if point.Point < req.Points {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "insufficient points"})
            return
        }

        point.Point -= req.Points
        if err := db.Save(&point).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update point record"})
            return
        }

        pointHistory := models.PointHistory{
            UserID:    userID,
            Point:     req.Points,
            Action:    "使用",
            Details:   req.Details,
            CreatedAt: time.Now(),
        }
        if err := db.Create(&pointHistory).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create point history record"})
            return
        }
    }
    c.JSON(http.StatusOK, gin.H{})
}


func GetHistory(c *gin.Context) {
    year, err := strconv.Atoi(c.Query("year"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
        return
    }

    month, err := strconv.Atoi(c.Query("month"))
    if err != nil || month < 1 || month > 12 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
        return
    }

    userId := c.Query("user_id")
    if userId == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }

    db, err := database.Connect()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    end := start.AddDate(0, 1, 0)
    var history []models.PointHistory
    db.Where("user_id = ? AND created_at >= ? AND created_at < ?", userId, start, end).Find(&history)

    for i, record := range history {
        formattedDate := record.CreatedAt.Format("2006-01-02")
        history[i].CreatedAtFormatted = formattedDate
    }
    c.JSON(http.StatusOK, history)
}

func GetPoint(c *gin.Context) {
    userId := c.Query("user_id")
    if userId == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }

    db, err := database.Connect()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    var point models.Point
    db.Where("user_id = ?", userId).First(&point)
    c.JSON(http.StatusOK, point)
}

func DeleteUser(c *gin.Context) {
    // リクエストからuser_idを取得
    userId := c.PostForm("user-id")

    // user_idが不正な場合、エラーレスポンスを返す
    if userId == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // データベースで対応するユーザーのdelete_flagをtrueに設定
    db, err := database.Connect()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    var user models.User
    if err := db.Where("id = ?", userId).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
        }
        return
    }

    user.DeleteFlag = true
    if err := db.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
        return
    }

    // 成功レスポンスを返す
    c.JSON(http.StatusOK, gin.H{"status": "User deleted", "message": "ユーザーが削除されました"})

}
