package models

type (
	// UserModel 用户模型
	UserModel struct {
		GormModel
		Username string `gorm:"type:varchar(128);not null;comment:;" json:"username"`
		Password string `gorm:"type:varchar(256);not null;comment:;" json:"password"`
		Nickname string `gorm:"type:varchar(128);not null;comment:;" json:"nickname"`
	}
)

// TableName 表名称
func (UserModel) TableName() string {
	return "users"
}
