package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dormsystem/db"
)

type BuildingOccupancyStat struct {
	BuildingID     uint    `json:"buildingID"`
	BuildingNo     string  `json:"buildingNo"`
	TotalCapacity  int     `json:"totalCapacity"`
	OccupiedBeds   int     `json:"occupiedBeds"`
	OccupancyRate  float64 `json:"occupancyRate"`
}

type BuildingPaymentSummary struct {
	BuildingID  uint    `json:"buildingID"`
	BuildingNo  string  `json:"buildingNo"`
	TotalAmount float64 `json:"totalAmount"`
}

func GetBuildingOccupancy(c *gin.Context) {
	var list []BuildingOccupancyStat
	if err := db.DB.Raw(`select building_id, building_no, total_capacity, occupied_beds, occupancy_rate from v_building_occupancy`).Scan(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询入住率统计失败"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func GetBuildingPaymentSummary(c *gin.Context) {
	var list []BuildingPaymentSummary
	if err := db.DB.Raw(`select building_id, building_no, total_amount from v_building_payment_summary`).Scan(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询收费统计失败"})
		return
	}
	c.JSON(http.StatusOK, list)
}

