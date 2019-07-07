package controllers

import (
	"fmt"
	"github.com/acidobinario/Angelhack-2019-Backend/database"
	"github.com/acidobinario/Angelhack-2019-Backend/login"
	"github.com/acidobinario/Angelhack-2019-Backend/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

/*
TODO:
	- Get Company Data
	- Get Company Balance
	- Get Company Employees (with its Balances, Salary, Full Name, ...)
	-
*/

func GetUser(context *gin.Context) {
	//TODO: implement me uwu
	JwtCtx := login.ExtractClaims(context)
	// identity claim is the username in the database
	username := JwtCtx[jwt.IdentityKey]
	userData := database.GetUser(username.(string))
	payload := models.UserInfoPayload{
		UserCommon: userData.UserCommon,
	}
	context.JSON(http.StatusOK, payload)
}

func AddEmployee(context *gin.Context) {
	JwtCtx := login.ExtractClaims(context)
	// identity claim is the username in the database
	username := JwtCtx[jwt.IdentityKey]
	thisCompany := database.GetUser(username.(string))
	var newEmployee struct {
		User   models.User `json:"user"`
		Salary float64     `json:"salary"`
	}
	//var newEmployee models.User
	err := context.ShouldBindJSON(&newEmployee)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"error binding json request": err})
		return
	}
	newEmployee.User.Access = "Usuario"
	hashedstuff, err := bcrypt.GenerateFromPassword([]byte(newEmployee.User.Password), 10)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"error": "could not create password"})
		return
	}
	newEmployee.User.Password = string(hashedstuff)
	database.CreateEmployee(newEmployee.User, newEmployee.Salary, thisCompany)
	context.JSON(http.StatusOK, gin.H{"status": "created!"})
}

func AddCompanyUser(context *gin.Context) {
	JwtCtx := login.ExtractClaims(context)
	// identity claim is the username in the database
	username := JwtCtx[jwt.IdentityKey]
	if username == "admin" {
		//TODO parse body to user struct
		var newUser models.User
		if err := context.BindJSON(&newUser); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error parsing JSON": err})
			return
		}
		newUser.Access = "Empresa"
		hashedstuff, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
		if err != nil {
			context.JSON(http.StatusBadGateway, gin.H{"error": "could not create password"})
			return
		}
		newUser.Password = string(hashedstuff)
		fmt.Println(newUser)

		database.CreateUser(newUser)
		context.JSON(http.StatusOK, gin.H{"status": "created!"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": "not admin"})
}

func GetEmployees(ctx *gin.Context) {
	thisAuthUser := getCurrentAuthUser(ctx)
	err, employees := database.GetEmployees(thisAuthUser.Username) //get the employeeeeet from this Company
	fmt.Println(employees)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "cannot get employees"})
		return
	}
	ctx.JSON(http.StatusOK, employees)
}

func getCurrentAuthUser(ctx *gin.Context) models.User {
	JwtCtx := login.ExtractClaims(ctx)
	// identity claim is the username in the database
	username := JwtCtx[jwt.IdentityKey]
	return database.GetUser(username.(string))
}
