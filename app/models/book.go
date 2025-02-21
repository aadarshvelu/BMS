package models

import "time"

// Book Model
type Book struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`	
	Title     string    `gorm:"not null" json:"title"`
	Author    string    `gorm:"not null" json:"author"`
	Year      int       `gorm:"not null" json:"year"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	IsActive  bool      `gorm:"default:true" json:"isActive"`
}

