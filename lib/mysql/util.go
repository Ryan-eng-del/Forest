package lib

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func WithContextAndTable (c *gin.Context, tableName string) func(db *gorm.DB) *gorm.DB{
  return func(db *gorm.DB) *gorm.DB {
    query := db.WithContext(c)
    return query.Table(tableName)
  }
}

func LogicalObjects() func(db *gorm.DB) *gorm.DB{
  return func(db *gorm.DB) *gorm.DB {
    return db.Where("is_deleted = ?", 0)
  }
}

func IDDesc() func(db *gorm.DB) *gorm.DB{
  return func(db *gorm.DB) *gorm.DB {
    return db.Order("id desc")
  }
}


func Paginate (pageNo, pageSize int) func(db *gorm.DB) *gorm.DB{
  return func(db *gorm.DB) *gorm.DB {
    if pageNo <= 0 {
      pageNo = 1
    }
    switch {
    case pageSize > 100:
      pageSize = 100
    case pageSize <= 0:
      pageSize = 10
    }

    offset := (pageNo - 1) * pageSize
    return db.Offset(offset).Limit(pageSize)
  }
}