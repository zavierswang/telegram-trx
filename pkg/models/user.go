package models

import "time"

type User struct {
	UserID    string    `gorm:"column:user_id;primaryKey;size:16"`
	Username  string    `gorm:"column:username;size:64;comment:用户名"`
	FirstName string    `gorm:"column:first_name;default:'';size:64;comment:用户别名"`
	IsAdmin   bool      `gorm:"column:is_admin;default:false;comment:是否为管理员"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime;comment:最新操作时间"`
}

func (u *User) TableName() string {
	return "tb_user"
}
