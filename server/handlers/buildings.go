package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"dormsystem/db"
	"dormsystem/models"
)

type BuildingRequest struct {
	BuildingNo string `json:"buildingNo"`
	FloorCount int    `json:"floorCount"`
	RoomCount  int    `json:"roomCount"`
	StartedAt  string `json:"startedAt"`
}

func ListBuildings(c *gin.Context) {
	var list []models.ApartmentBuilding
	if applyPagination(c, db.DB.Model(&models.ApartmentBuilding{}), &list) {
		return
	}
}

func CreateBuilding(c *gin.Context) {
	var req BuildingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据不合法"})
		return
	}
	startedAt, err := time.Parse("2006-01-02", req.StartedAt)
	if err != nil {
		startedAt, err = time.Parse(time.RFC3339, req.StartedAt)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "启用时间格式应为YYYY-MM-DD"})
		return
	}
	b := models.ApartmentBuilding{
		BuildingNo: req.BuildingNo,
		FloorCount: req.FloorCount,
		RoomCount:  req.RoomCount,
		StartedAt:  startedAt,
	}
	if err := db.DB.Create(&b).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, b)
}

func UpdateBuilding(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	var b models.ApartmentBuilding
	if err := db.DB.First(&b, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}
	var req BuildingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据不合法"})
		return
	}
	startedAt, err := time.Parse("2006-01-02", req.StartedAt)
	if err != nil {
		startedAt, err = time.Parse(time.RFC3339, req.StartedAt)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "启用时间格式应为YYYY-MM-DD"})
		return
	}
	b.BuildingNo = req.BuildingNo
	b.FloorCount = req.FloorCount
	b.RoomCount = req.RoomCount
	b.StartedAt = startedAt
	if err := db.DB.Save(&b).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, b)
}

func DeleteBuilding(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	if err := db.DB.Delete(&models.ApartmentBuilding{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
