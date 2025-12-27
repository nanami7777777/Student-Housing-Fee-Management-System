package models

import "time"

type ApartmentBuilding struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	BuildingNo string    `gorm:"uniqueIndex;size:20;not null" json:"buildingNo"`
	FloorCount int       `gorm:"not null" json:"floorCount"`
	RoomCount  int       `gorm:"not null" json:"roomCount"`
	StartedAt  time.Time `gorm:"not null;type:date" json:"startedAt"`
	Rooms      []DormRoom `gorm:"foreignKey:BuildingID;references:ID" json:"-"`
}

type DormRoom struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	RoomNo     string  `gorm:"size:20;not null" json:"roomNo"`
	Capacity   int     `gorm:"not null" json:"capacity"`
	Fee        float64 `json:"fee"`
	Phone      string  `gorm:"size:20" json:"phone"`
	BuildingID uint    `gorm:"not null;index" json:"buildingID"`
	Building   ApartmentBuilding `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Students   []Student         `gorm:"foreignKey:RoomID;references:ID" json:"-"`
	Payments   []Payment         `gorm:"foreignKey:RoomID;references:ID" json:"-"`
}

type Student struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	StudentNo  string `gorm:"uniqueIndex;size:20;not null" json:"studentNo"`
	Name       string `gorm:"size:50;not null" json:"name"`
	Gender     string `gorm:"size:10" json:"gender"`
	Ethnicity  string `gorm:"size:20" json:"ethnicity"`
	Major      string `gorm:"size:100" json:"major"`
	ClassName  string `gorm:"size:50" json:"className"`
	Phone      string `gorm:"size:20" json:"phone"`
	BuildingID uint   `gorm:"index" json:"buildingID"`
	RoomID     uint   `gorm:"index" json:"roomID"`
	Building   ApartmentBuilding `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Room       DormRoom          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Payments   []Payment         `json:"-"`
}

type Payment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PaymentNo   string    `gorm:"uniqueIndex;size:50;not null" json:"paymentNo"`
	BuildingID  uint      `gorm:"not null;index" json:"buildingID"`
	RoomID      uint      `gorm:"not null;index" json:"roomID"`
	StudentID   uint      `gorm:"index" json:"studentID"`
	PaidAt      time.Time `gorm:"not null;type:date" json:"paidAt"`
	PaymentType string    `gorm:"size:50;not null" json:"paymentType"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Building    ApartmentBuilding `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Room        DormRoom          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Student     Student           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Name         string    `gorm:"size:50;not null" json:"name"`
	PasswordHash string    `gorm:"size:200;not null" json:"-"`
	Role         string    `gorm:"size:20;not null" json:"role"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
