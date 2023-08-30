package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	Address string  `gorm:"column:address;uniqueIndex;size:64"` //预支地址不能重复
	Balance float64 `gorm:"column:balance;default:0.0"`         //兑换总额
	Advance float64 `gorm:"column:advance;default:0.0"`         //预金额（记录次数）
	Count   int     `gorm:"column:count;default:0"`             //可预支次数
}

func (a *Address) TableName() string {
	return "tb_address"
}
