package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"dormsystem/db"
	"dormsystem/models"
)

type PaymentRequest struct {
	PaymentNo   string  `json:"paymentNo"`
	BuildingID  uint    `json:"buildingID"`
	RoomID      uint    `json:"roomID"`
	StudentID   uint    `json:"studentID"`
	PaidAt      string  `json:"paidAt"`
	PaymentType string  `json:"paymentType"`
	Amount      float64 `json:"amount"`
}

func ListPayments(c *gin.Context) {
	var list []models.Payment
	query := db.DB.Model(&models.Payment{})
	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where(
			db.DB.
				Where("payment_no = ?", keyword).
				Or("payment_type = ?", keyword),
		)
	}
	buildingIDStr := c.Query("buildingID")
	if buildingIDStr != "" {
		if id, err := strconv.Atoi(buildingIDStr); err == nil && id > 0 {
			query = query.Where("building_id = ?", id)
		}
	}
	roomIDStr := c.Query("roomID")
	if roomIDStr != "" {
		if id, err := strconv.Atoi(roomIDStr); err == nil && id > 0 {
			query = query.Where("room_id = ?", id)
		}
	}
	studentIDStr := c.Query("studentID")
	if studentIDStr != "" {
		if id, err := strconv.Atoi(studentIDStr); err == nil && id > 0 {
			query = query.Where("student_id = ?", id)
		}
	}
	if applyPagination(c, query, &list) {
		return
	}
}

func CreatePayment(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据不合法"})
		return
	}
	if req.BuildingID == 0 || req.RoomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "公寓号和寝室号不能为空"})
		return
	}
	var b models.ApartmentBuilding
	if err := db.DB.First(&b, req.BuildingID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "公寓不存在"})
		return
	}
	var r models.DormRoom
	if err := db.DB.First(&r, req.RoomID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "寝室不存在"})
		return
	}
	if r.BuildingID != req.BuildingID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "寝室不属于该公寓"})
		return
	}
	if req.StudentID != 0 {
		var s models.Student
		if err := db.DB.First(&s, req.StudentID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "学生不存在"})
			return
		}
		if s.RoomID != req.RoomID || s.BuildingID != req.BuildingID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "学生不在指定公寓寝室中"})
			return
		}
	}
	paidAt, err := time.Parse("2006-01-02", req.PaidAt)
	if err != nil {
		paidAt, err = time.Parse(time.RFC3339, req.PaidAt)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "交费日期格式应为YYYY-MM-DD"})
		return
	}
	p := models.Payment{
		PaymentNo:   req.PaymentNo,
		BuildingID:  req.BuildingID,
		RoomID:      req.RoomID,
		StudentID:   req.StudentID,
		PaidAt:      paidAt,
		PaymentType: req.PaymentType,
		Amount:      req.Amount,
	}
	if err := db.DB.Create(&p).Error; err != nil {
		respondDBError(c, err, "创建失败")
		return
	}
	c.JSON(http.StatusOK, p)
}

func UpdatePayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	var p models.Payment
	if err := db.DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据不合法"})
		return
	}
	if req.BuildingID == 0 || req.RoomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "公寓号和寝室号不能为空"})
		return
	}
	var b models.ApartmentBuilding
	if err := db.DB.First(&b, req.BuildingID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "公寓不存在"})
		return
	}
	var r models.DormRoom
	if err := db.DB.First(&r, req.RoomID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "寝室不存在"})
		return
	}
	if r.BuildingID != req.BuildingID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "寝室不属于该公寓"})
		return
	}
	if req.StudentID != 0 {
		var s models.Student
		if err := db.DB.First(&s, req.StudentID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "学生不存在"})
			return
		}
		if s.RoomID != req.RoomID || s.BuildingID != req.BuildingID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "学生不在指定公寓寝室中"})
			return
		}
	}
	paidAt, err := time.Parse("2006-01-02", req.PaidAt)
	if err != nil {
		paidAt, err = time.Parse(time.RFC3339, req.PaidAt)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "交费日期格式应为YYYY-MM-DD"})
		return
	}
	p.PaymentNo = req.PaymentNo
	p.BuildingID = req.BuildingID
	p.RoomID = req.RoomID
	p.StudentID = req.StudentID
	p.PaidAt = paidAt
	p.PaymentType = req.PaymentType
	p.Amount = req.Amount
	if err := db.DB.Save(&p).Error; err != nil {
		respondDBError(c, err, "更新失败")
		return
	}
	c.JSON(http.StatusOK, p)
}

func DeletePayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	if err := db.DB.Delete(&models.Payment{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
