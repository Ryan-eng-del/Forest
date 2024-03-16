package lib

import "gorm.io/gorm"

func LogicalObjects(db *gorm.DB) *gorm.DB {
  return db.Where("is_deleted = ?", 0)
}