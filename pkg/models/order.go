package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	TxID        string  `gorm:"column:tx_id;uniqueIndex;size:64"`
	FromAddress string  `gorm:"column:from_address;not null"`
	ToAddress   string  `gorm:"column:to_address;not null"`
	Balance     float64 `gorm:"column:balance;default:0.0"`
	Amount      float64 `gorm:"column:amount;default:0.0"`
	Finished    bool    `gorm:"column:finished;default:false"`
	Status      int     `gorm:"column:status"`
}

func (t *Order) TableName() string {
	return "tb_order"
}
