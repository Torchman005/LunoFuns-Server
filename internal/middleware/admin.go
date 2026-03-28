package middleware

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    
    "LunoFuns-Server/internal/models"
)

func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "未认证",
            })
            c.Abort()
            return
        }
        
        if role.(int8) != models.UserRoleAdmin {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "权限不足",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}