package main

import (
	"fmt"
	"github.com/acidobinario/Angelhack-2019-Backend/controllers"
	"github.com/acidobinario/Angelhack-2019-Backend/login"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

/*
	/api/
		login POST
		users/ GET(get user stuff)
		auth/ (JWT)
			user


*/
func main() {

	auth, err := login.GetJWT()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("owowowowowow")
	g := gin.Default()
	g.LoadHTMLGlob("./public/templates/**/*")

	api := g.Group("/api")
	api.POST("/login", auth.LoginHandler)
	api.GET("/refresh_token", auth.RefreshHandler)
	api.POST("/verify", controllers.VerifyQRCode)     // Validate QR code (Agent type)
	authEndpoint := api.Group("/auth").Use(auth.MiddlewareFunc())
	authEndpoint.GET("/test", func(context *gin.Context) {
		fmt.Println(login.ExtractClaims(context))
		context.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authEndpoint.GET("/user", controllers.GetUser)               //retrieve the user info
	authEndpoint.GET("/pay", controllers.PayEmployees)           //endpoint to pay to its employees
	authEndpoint.POST("/addEmployee", controllers.AddEmployee)   // only admin can add Company type User
	authEndpoint.POST("/addCompany", controllers.AddCompanyUser) //Companies can add Employee type Users
	authEndpoint.GET("/employees", controllers.GetEmployees)

	g.Static("/assets", "./public/assets/")
	g.GET("/login", func(context *gin.Context) { // login controler
		context.HTML(http.StatusOK, "login.html", nil)
	})
	g.GET("/index", func(context *gin.Context) { // login controler
		context.HTML(http.StatusOK, "index.html", nil)
	})
	g.GET("/employees", func(context *gin.Context) { // login controler
		context.HTML(http.StatusOK, "employees.html", nil)
	})
	g.GET("/transactions", func(context *gin.Context) { // login controler
		context.HTML(http.StatusOK, "transactions.html", nil)
	})
	g.GET("/verify", func(context *gin.Context) { // login controler
		context.HTML(http.StatusOK, "verify.gohtml", nil)
	})
	if err := g.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
