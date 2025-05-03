package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"oj/helper"
)

func AuthUserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		userClaim, err := helper.AnalyseToken(auth)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized Authorization",
			})
			return
		}
		if userClaim == nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized User",
			})
			return
		}
		fmt.Println(userClaim)
		c.Set("user_claims", userClaim)
		c.Next()
	}
}
