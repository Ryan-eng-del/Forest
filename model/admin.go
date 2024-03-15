package model



type Admin struct {
	AbstractModel
	UserName string `json:"user_name" gorm:"type:varchar(255);not null;comment:用户名" `	
	Salt string `json:"salt" gorm:"type:varchar(50);not null;comment:盐值" `
	Password string `json:"password" gorm:"type:varchar(255);comment:密码;not null"`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}
