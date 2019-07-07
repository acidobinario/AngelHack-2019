package login

import (
	"github.com/acidobinario/Angelhack-2019-Backend/database"
	"github.com/acidobinario/Angelhack-2019-Backend/models"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"time"
)

func GetJWT() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "AngelHack Bank",
		Key:        []byte("angelHack2018"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.Login); ok {
				return jwt.MapClaims{
					jwt.IdentityKey: v.Username,
					"access": v.Access,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(context *gin.Context) interface{} {
			claims := jwt.ExtractClaims(context)
			return &models.User{
				Login: models.Login{
					Username: claims["identity"].(string),
				},
			}
		},
		Authenticator: func(c *gin.Context) (i interface{}, e error) {
			var loginVals models.Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if database.CheckUser(userID, password) == true {
				user := database.GetUser(userID)
				return &models.Login{
					Username: userID,
					Access:user.Access,
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			//userExists = connect.UsernameExists()
			//if v, ok := data.(*models.User); ok && v.Username == "admin" {
			if v, ok := data.(*models.User); ok && database.UsernameExists(v.Username) {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	})
}

func ExtractClaims(ctx *gin.Context) jwt.MapClaims{
	return jwt.ExtractClaims(ctx)
}
