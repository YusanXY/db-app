package middleware

import (
	"dbapp/internal/errors"
	"dbapp/pkg/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errors.HandleError(c, errors.NewUnauthorizedError("未提供认证信息"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			errors.HandleError(c, errors.NewUnauthorizedError("认证格式错误"))
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := utils.ParseJWT(token)
		if err != nil {
			errors.HandleError(c, errors.NewUnauthorizedError("Token无效或已过期"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

