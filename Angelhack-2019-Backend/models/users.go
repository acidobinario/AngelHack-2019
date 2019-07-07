package models

import "github.com/jinzhu/gorm"

type User struct {
	//tabla users
	gorm.Model
	Login
	UserCommon
}

type UserCommon struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Rut     int    `json:"rut"` //this shouldn't be rut lmao, more like DocumentNumber or something like that
	Address string `json:"address"`
	Email   string `json:"email"`
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Access   string // access can be either  1: Empresa, 2: Usuario or 3: Agente
}

type Employee struct {
	gorm.Model
	CompanyId uint
	UserId    uint
	Balance   float64
	Code      string
	IsUsed    bool
	Salary    float64
}

/*
	this is the payload response through the /user endpoint when requesting the info about the user
*/
type UserInfoPayload struct {
	UserCommon
}


type EmployeeList struct {
	User
	Employee
}

/*
Login with contact n stuff
Table users:
Id, Username, Password, FullName, Type, Email,...
1, TeleWorm, 11235, TeleWorm Inc, Company, us@TeleWorm.us, ..
2, Calvin, 12378, Calvin M. Murray, User, CalvinMMurray@teleworm.us, ...
3, David, 81wdiu, David V. Harris, User, DavidVHarris@teleworm.us, ...
4, dayrep, g8f21, dayrep Inc, Company, JohnnieCBrown@dayrep.com, ...
5, Amons1947, vooYai9iy9Ie, Johnnie C. Brown, User, JohnnieCBrown@dayrep.com
*/

/*
Info de correlacion company employee
table Employees 2:
CompanyId, UserId, Balance
1, 2, 0
1, 3, 0
4, 5, 0
*/
