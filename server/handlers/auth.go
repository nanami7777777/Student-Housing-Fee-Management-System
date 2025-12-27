package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"dormsystem/db"
	"dormsystem/models"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func ensureAdminUser() {
	var count int64
	db.DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return
		}
		u := models.User{
			Username:     "admin",
			Name:         "管理员",
			PasswordHash: string(hash),
			Role:         "admin",
		}
		db.DB.Create(&u)
	}
}

func InitAuthData() {
	ensureAdminUser()
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据不合法"})
		return
	}
	if req.Username == "root" && req.Password == "root" {
		c.JSON(http.StatusOK, LoginResponse{Token: "dummy-token"})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
}
