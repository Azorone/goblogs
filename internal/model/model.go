package model

import (
	"time"
)

type User struct {
	ID       uint   `gorm:"primarykey"`
	Username string `gorm:"unique" json:"Username,omitempty"`
	Password string `json:"Password,omitempty"`
}

type Blog struct {
	ID         uint       `gorm:"primarykey"`
	CreatedAt  *time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt  *time.Time `json:"UpdatedAt,omitempty"`
	DeletedAt  *time.Time `gorm:"index" gorm:"default:null" json:"DeletedAt,omitempty"`
	Title      string     `json:"Title,omitempty"`
	Content    string     `json:"Content,omitempty"`
	CategoryID uint       `json:"Category_id,omitempty"`
	UserID     uint       `json:"User_id,omitempty"`
	Deleted    bool       `gorm:"Default:false" json:"deleted,omitempty"`
	Des        string     `json:"Des,omitempty"`
	IsPush     bool       `gorm:"default:false" json:"IsPush,omitempty"`
}

type Category struct {
	ID        uint       `gorm:"primarykey"`
	CreatedAt *time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt *time.Time `json:"UpdatedAt,omitempty"`
	DeletedAt *time.Time `gorm:"index" gorm:"default:null" json:"DeletedAt,omitempty"`
	Name      string     `gorm:"type:varchar(255);unique"`
}

type Image struct {
	ID   string `gorm:"type:varchar(255);primary_key"`
	Path string `gorm:"type:varchar(255)"`
}
type AccessLog struct {
	ID           uint          `gorm:"primarykey"`
	AccessTime   time.Time     `json:"access_time"`
	IP           string        `json:"ip"`
	Method       string        `json:"method"`
	Path         string        `json:"path"`
	ResponseTime time.Duration `json:"response_time"`
	ErrorMsg     string        `json:"error_msg"`
}
