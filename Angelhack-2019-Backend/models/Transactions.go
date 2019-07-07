package models

import "github.com/jinzhu/gorm"
/*
transactions logs
table Transactions:
ID, From, To, Amount
*/
type Transaction struct {
	gorm.Model
	From   uint
	To     uint
	Amount float64
}

