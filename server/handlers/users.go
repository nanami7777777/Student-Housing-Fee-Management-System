package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"dormsystem/db"
	"dormsystem/models"
)

type UserRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func ListUsers(c *gin.Context) {
	var list []models.User
	query := db.DB.Model(&models.User{})
	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where(
			db.DB.
				Where("username = ?", keyword).
				Or("name = ?", keyword),
		)
	}
	role := c.Query("role")
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if applyPagination(c, query, &list) {
		return
	}
}

func CreateUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据不合法"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码处理失败"})
		return
	}
	u := models.User{
		Username:     req.Username,
		Name:         req.Name,
		PasswordHash: string(hash),
		Role:         req.Role,
	}
	if err := db.DB.Create(&u).Error; err != nil {
		respondDBError(c, err, "创建失败")
		return
	}
	c.JSON(http.StatusOK, u)
}

func UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	var u models.User
	if err := db.DB.First(&u, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据不合法"})
		return
	}
	u.Username = req.Username
	u.Name = req.Name
	u.Role = req.Role
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码处理失败"})
			return
		}
		u.PasswordHash = string(hash)
	}
	if err := db.DB.Save(&u).Error; err != nil {
		respondDBError(c, err, "更新失败")
		return
	}
	c.JSON(http.StatusOK, u)
}

func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	if err := db.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
