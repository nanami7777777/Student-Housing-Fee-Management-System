package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	pgxpgconn "github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"dormsystem/db"
	"dormsystem/models"
)

func paginateParams(c *gin.Context) (int, int, bool) {
	pageStr := c.Query("page")
	sizeStr := c.Query("pageSize")
	if pageStr == "" && sizeStr == "" {
		return 0, 0, false
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	offset := (page - 1) * size
	return size, offset, true
}

func applyPagination(c *gin.Context, query *gorm.DB, out interface{}) bool {
	limit, offset, usePagination := paginateParams(c)
	if !usePagination {
		if err := query.Find(out).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
			return false
		}
		c.JSON(http.StatusOK, out)
		return false
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return false
	}
	if err := query.Limit(limit).Offset(offset).Find(out).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return false
	}
	c.JSON(http.StatusOK, gin.H{
		"items": out,
		"total": total,
	})
	return true
}

func respondDBError(c *gin.Context, err error, defaultMsg string) {
	msg := err.Error()
	constraint := ""

	var pgErr *pgxpgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.ConstraintName != "" {
			constraint = pgErr.ConstraintName
		}
		if msg == "" {
			msg = pgErr.Message
		}
	}

	if constraint == "" {
		switch {
		case strings.Contains(msg, "chk_building_floor_count"):
			constraint = "chk_building_floor_count"
		case strings.Contains(msg, "chk_building_room_count"):
			constraint = "chk_building_room_count"
		case strings.Contains(msg, "chk_room_capacity"):
			constraint = "chk_room_capacity"
		case strings.Contains(msg, "chk_room_fee_nonnegative"):
			constraint = "chk_room_fee_nonnegative"
		case strings.Contains(msg, "chk_student_gender"):
			constraint = "chk_student_gender"
		case strings.Contains(msg, "chk_payment_amount"):
			constraint = "chk_payment_amount"
		case strings.Contains(msg, "chk_payment_type"):
			constraint = "chk_payment_type"
		}
	}

	switch constraint {
	case "chk_building_floor_count":
		c.JSON(http.StatusBadRequest, gin.H{"error": "楼层数必须大于0"})
		return
	case "chk_building_room_count":
		c.JSON(http.StatusBadRequest, gin.H{"error": "房间数不能为负数"})
		return
	case "chk_room_capacity":
		c.JSON(http.StatusBadRequest, gin.H{"error": "寝室容量必须大于0"})
		return
	case "chk_room_fee_nonnegative":
		c.JSON(http.StatusBadRequest, gin.H{"error": "寝室费用不能为负数"})
		return
	case "chk_student_gender":
		c.JSON(http.StatusBadRequest, gin.H{"error": "性别只能是男或女"})
		return
	case "chk_payment_amount":
		c.JSON(http.StatusBadRequest, gin.H{"error": "金额必须大于0"})
		return
	case "chk_payment_type":
		c.JSON(http.StatusBadRequest, gin.H{"error": "收费类型必须是住宿费、水电费或押金"})
		return
	default:
		if strings.Contains(msg, "寝室不存在") || strings.Contains(msg, "寝室人数已满") {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": defaultMsg})
}

func ListStudents(c *gin.Context) {
	var students []models.Student
	query := db.DB.Model(&models.Student{})
	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword != "" {
		query = query.Where(
			db.DB.
				Where("student_no = ?", keyword).
				Or("name = ?", keyword).
				Or("class_name = ?", keyword),
		)
	}
	roomNo := strings.TrimSpace(c.Query("roomNo"))
	if roomNo != "" {
		query = query.Joins("JOIN dorm_rooms ON dorm_rooms.id = students.room_id").
			Where("dorm_rooms.room_no = ?", roomNo)
	}
	if applyPagination(c, query, &students) {
		return
	}
}

func CreateStudent(c *gin.Context) {
	var s models.Student
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据不合法"})
		return
	}
	if s.BuildingID == 0 || s.RoomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "公寓号和寝室号不能为空"})
		return
	}
	var b models.ApartmentBuilding
	if err := db.DB.First(&b, s.BuildingID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "公寓不存在"})
		return
	}
	var r models.DormRoom
	if err := db.DB.First(&r, s.RoomID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "寝室不存在"})
		return
	}
	if r.BuildingID != s.BuildingID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "寝室不属于该公寓"})
		return
	}
	if err := db.DB.Create(&s).Error; err != nil {
		respondDBError(c, err, "创建失败")
		return
	}
	c.JSON(http.StatusOK, s)
}

func UpdateStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	var s models.Student
	if err := db.DB.First(&s, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}
	var req models.Student
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
	s.StudentNo = req.StudentNo
	s.Name = req.Name
	s.Gender = req.Gender
	s.Ethnicity = req.Ethnicity
	s.Major = req.Major
	s.ClassName = req.ClassName
	s.Phone = req.Phone
	s.BuildingID = req.BuildingID
	s.RoomID = req.RoomID
	if err := db.DB.Save(&s).Error; err != nil {
		respondDBError(c, err, "更新失败")
		return
	}
	c.JSON(http.StatusOK, s)
}

func DeleteStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	if err := db.DB.Delete(&models.Student{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
