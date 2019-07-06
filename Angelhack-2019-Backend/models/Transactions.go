package models

import "github.com/jinzhu/gorm"

type Transaction struct {
	gorm.Model
	From   int
	To     int
	Amount float64
}
