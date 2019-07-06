package models

import "github.com/jinzhu/gorm"

type User struct { //tabla users
	gorm.Model
	Login
}

type UserCommon struct {
	Name string
	Phone string
	Rut int
	Address string
}

type Login struct {
	Username string
	Password string
	Access string // access can be either  1: Empresa, 2: Usuario or 3: Agente
}

type UserEmpresa struct {
	gorm.Model
	UserCommon
}


