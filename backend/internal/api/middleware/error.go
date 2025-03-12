package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/yourusername/crypto-monitor/internal/errors"
)

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            switch e := err.Err.(type) {
            case *errors.ValidationError:
                c.JSON(http.StatusBadRequest, gin.H{
                    "error": e.Error(),
                })
            case *errors.AuthenticationError:
                c.JSON(http.StatusUnauthorized, gin.H{
                    "error": e.Error(),
                })
            case *errors.AuthorizationError:
                c.JSON(http.StatusForbidden, gin.H{
                    "error": e.Error(),
                })
            default:
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "Internal server error",
                })
            }
        }
    }
}