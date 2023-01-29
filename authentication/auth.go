package authentication

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	UsernameKey = "username"
	EmailKey    = "email"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := strings.Split(context.Request.Header.Get("Authorization"), "Bearer ")
		if len(tokenString) != 2 {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		token, err := ValidateToken(tokenString[1])
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.Set(UsernameKey, token.Username)
		context.Set(EmailKey, token.Email)
		context.Next()
	}
}
