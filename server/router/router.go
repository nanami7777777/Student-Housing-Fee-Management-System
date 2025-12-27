package router

import (
	"github.com/gin-gonic/gin"

	"dormsystem/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	r.POST("/api/login", handlers.Login)
	api := r.Group("/api")
	api.GET("/stats/building-occupancy", handlers.GetBuildingOccupancy)
	api.GET("/stats/building-payments", handlers.GetBuildingPaymentSummary)
	api.GET("/buildings", handlers.ListBuildings)
	api.POST("/buildings", handlers.CreateBuilding)
	api.PUT("/buildings/:id", handlers.UpdateBuilding)
	api.DELETE("/buildings/:id", handlers.DeleteBuilding)
	api.GET("/rooms", handlers.ListRooms)
	api.POST("/rooms", handlers.CreateRoom)
	api.PUT("/rooms/:id", handlers.UpdateRoom)
	api.DELETE("/rooms/:id", handlers.DeleteRoom)
	api.GET("/students", handlers.ListStudents)
	api.POST("/students", handlers.CreateStudent)
	api.PUT("/students/:id", handlers.UpdateStudent)
	api.DELETE("/students/:id", handlers.DeleteStudent)
	api.GET("/payments", handlers.ListPayments)
	api.POST("/payments", handlers.CreatePayment)
	api.PUT("/payments/:id", handlers.UpdatePayment)
	api.DELETE("/payments/:id", handlers.DeletePayment)
	api.GET("/users", handlers.ListUsers)
	api.POST("/users", handlers.CreateUser)
	api.PUT("/users/:id", handlers.UpdateUser)
	api.DELETE("/users/:id", handlers.DeleteUser)
	return r
}
