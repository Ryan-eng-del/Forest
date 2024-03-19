package model

import (
	"errors"
	adminDto "go-gateway/dto/admin"
	lib "go-gateway/lib/func"

	mysqlLib "go-gateway/lib/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type Admin struct {
	AbstractModel
	UserName string `json:"user_name" gorm:"type:varchar(255);not null;comment:用户名" `	
	Salt string `json:"salt" gorm:"type:varchar(50);not null;comment:盐值" `
	Password string `json:"password" gorm:"type:varchar(255);comment:密码;not null"`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}


func (t *Admin) LoginCheck(c *gin.Context, db *gorm.DB, params adminDto.AdminLoginInput ) (*Admin, error) {

	admin, err := t.FindByName(c, db, params.Username)
	if err != nil {
		return nil, errors.New("用户信息不存在")
	}

	saltPassword := lib.GenSaltPassword(admin.Salt, params.Password)

	if admin.Password != saltPassword {
		return nil, errors.New("密码输入不正确")
	}
	return admin, nil
}


func (t *Admin) FindByName(c *gin.Context, db *gorm.DB, search string) (*Admin, error) {
	out := &Admin{}
	err := db.Scopes(mysqlLib.WithContextAndTable(c, t.TableName()), mysqlLib.LogicalObjects()).WithContext(c).Where("user_name = ?", search).First(&out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *Admin) FindById(c *gin.Context, db *gorm.DB, id uint) (*Admin, error) {
	out := &Admin{}
	err := db.Scopes(mysqlLib.WithContextAndTable(c, t.TableName()), mysqlLib.LogicalObjects()).WithContext(c).First(out, id).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *Admin) Save(c *gin.Context, db *gorm.DB) error {
	return db.WithContext(c).Save(t).Error
}