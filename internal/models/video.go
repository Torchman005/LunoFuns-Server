package models

import (
	"time"
)

type Video struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64     `gorm:"not null" json:"user_id"`
	Title         string     `gorm:"type:varchar(256);not null" json:"title"`
	Description   string     `gorm:"type:text" json:"description"`
	CoverURL      string     `gorm:"type:varchar(512);not null" json:"cover_url"`
	VideoURL      string     `gorm:"type:varchar(512);not null" json:"video_url"`
	Duration      uint       `gorm:"not null" json:"duration"`
	Size          uint64     `gorm:"not null" json:"size"`
	Format        string     `gorm:"type:varchar(16);not null" json:"format"`
	Width         uint16     `gorm:"default:0" json:"width"`
	Height        uint16     `gorm:"default:0" json:"height"`
	ViewCount     uint64     `gorm:"default:0" json:"view_count"`
	DanmakuCount  uint       `gorm:"default:0" json:"danmaku_count"`
	LikeCount     uint       `gorm:"default:0" json:"like_count"`
	FavoriteCount uint       `gorm:"default:0" json:"favorite_count"`
	CoinCount     uint       `gorm:"default:0" json:"coin_count"`
	Status        int8       `gorm:"type:tinyint;default:1" json:"status"`
	CategoryID    uint       `gorm:"default:0" json:"category_id"`
	Tags          string     `gorm:"type:varchar(512);default:''" json:"tags"`
	PublishedAt   *time.Time `json:"published_at"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Video) TableName() string {
	return "video"
}