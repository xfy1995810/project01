package models

import (
	"gorm.io/gorm"
	"time"
)

type Common struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedID uint           `json:"created_id"`
	CreatedBy CommonSysUser  `gorm:"foreignKey:CreatedID;default:NULL;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;" json:"created_by"`
	CreatedAt time.Time      `json:"created_at" example:"创建时间"`
	UpdatedID uint           `json:"updated_id"`
	UpdatedBy CommonSysUser  `gorm:"foreignKey:UpdatedID;default:NULL;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID;" json:"updated_by"`
	UpdatedAt time.Time      `json:"updated_at" example:"更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;" json:"-"`
}

func HandleCommonUserPreload(db *gorm.DB) *gorm.DB {
	return db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB { return db.Select("ID", "Username", "ChineseName") }).
		Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB { return db.Select("ID", "Username", "ChineseName") })
}
