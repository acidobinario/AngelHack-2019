package login

import (
	"github.com/HappyVolt/Monitoreo-freezers-Backend/connect"
	"github.com/HappyVolt/Monitoreo-freezers-Backend/structures"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"time"
)

func GetJWT() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "Monitoreo freezers",
		Key:        []byte("desarrollo2482"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*structures.Login); ok {
				return jwt.MapClaims{
					jwt.IdentityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(context *gin.Context) interface{} {
			claims := jwt.ExtractClaims(context)
			return &structures.User{
				Login: structures.Login{
					Username: claims["identity"].(string),
				},
			}
		},
		Authenticator: func(c *gin.Context) (i interface{}, e error) {
			var loginVals structures.Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if connect.CheckUser(userID, password) == true {
				return &structures.Login{
					Username: userID,
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			//userExists = connect.UsernameExists()
			//if v, ok := data.(*structures.User); ok && v.Username == "admin" {
			if v, ok := data.(*structures.User); ok && connect.UsernameExists(v.Username) {
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
