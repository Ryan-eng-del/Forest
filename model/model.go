package model

import (
	"database/sql"
	"time"
)


type AbstractModel struct {
	ID        uint `gorm:"primarykey;comment:自增主键"`
	CreatedAt time.Time  `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
	DeletedAt  sql.NullTime `gorm:"comment:删除时间"`
	IsDeleted bool  `gorm:"comment:是否删除"`
}

