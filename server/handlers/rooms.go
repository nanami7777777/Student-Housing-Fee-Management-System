package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"dormsystem/db"
	"dormsystem/models"
)

func ListRooms(c *gin.Context) {
	var list []models.DormRoom
	query := db.DB.Model(&models.DormRoom{})
	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("room_no = ?", keyword)
	}
	buildingIDStr := c.Query("buildingID")
	if buildingIDStr != "" {
		if id, err := strconv.Atoi(buildingIDStr); err == nil && id > 0 {
			query = query.Where("building_id = ?", id)
		}
	}
	if applyPagination(c, query, &list) {
		return
	}
}

func CreateRoom(c *gin.Context) {
	var r models.DormRoom
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据不合法"})
		return
	}
	if r.BuildingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "所属公寓不能为空"})
		return
	}
	var b models.ApartmentBuilding
	if err := db.DB.First(&b, r.BuildingID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "所属公寓不存在"})
		return
	}
	if err := db.DB.Create(&r).Error; err != nil {
		respondDBError(c, err, "创建失败")
		return
	}
	c.JSON(http.StatusOK, r)
}

func UpdateRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	var r models.DormRoom
	if err := db.DB.First(&r, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}
	var req models.DormRoom
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据不合法"})
		return
	}
	if req.BuildingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "所属公寓不能为空"})
		return
	}
	var b models.ApartmentBuilding
	if err := db.DB.First(&b, req.BuildingID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "所属公寓不存在"})
		return
	}
	r.RoomNo = req.RoomNo
	r.Capacity = req.Capacity
	r.Fee = req.Fee
	r.Phone = req.Phone
	r.BuildingID = req.BuildingID
	if err := db.DB.Save(&r).Error; err != nil {
		respondDBError(c, err, "更新失败")
		return
	}
	c.JSON(http.StatusOK, r)
}

func DeleteRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	if err := db.DB.Delete(&models.DormRoom{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
