package controllers

import (
	"fmt"
	"github.com/acidobinario/Angelhack-2019-Backend/database"
	"github.com/acidobinario/Angelhack-2019-Backend/login"
	"github.com/acidobinario/Angelhack-2019-Backend/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

var fee = 15.0

func GetTransactionsByUserId(userId int) []models.Transaction {
	return nil
}

func VerifyQRCode(context *gin.Context) {
	//TODO: Implement Me!!!?")
	//get the Code
	var code struct{
		Code string `json:"code"`
	}
	if err := context.ShouldBindJSON(&code); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error!": err})
		return
	}
	fmt.Println(code)
	aaaa, err := database.VerifyQr(code.Code)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error!": err})
		return
	}
	database.SetCodeUsed(aaaa)
	context.JSON(http.StatusOK, aaaa)
	return
}

//func getQRCode
func PayEmployees(context *gin.Context) {
	JwtCtx := login.ExtractClaims(context)
	// identity claim is the username in the database
	access := JwtCtx["access"]
	if access.(string) == "Empresa" {
		//if the access is correct, then pay the Employeeeeeeeee
		thisCompany := database.GetUser(JwtCtx[jwt.IdentityKey].(string))
		err, employees := database.GetEmployees(JwtCtx[jwt.IdentityKey].(string)) //get the employeeeeet from this Company
		if err != nil {
			context.JSON(http.StatusBadGateway, gin.H{"error": "cannot get employees"})
			return
		}
		//for each one of them, send a transaction and generate the code for retirar the money
		for _, v := range employees {
			sendTransaction(thisCompany.ID, v.UserId, v.Salary) // {"code": "{codigo}"} {"status": "ok", "value":400000.02}

		}
	}
}

func sendTransaction(from uint, to uint, amount float64) {
	database.CreateTransaction(from, to, amount)
	database.GenerateCode(from)
	sendQrCode(to)
}

func sendQrCode(u uint) {
	//thisEmployee
}
