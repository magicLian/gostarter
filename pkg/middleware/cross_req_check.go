package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/magicLian/gostarter/pkg/setting"
)

func CrossRequestCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		passed := false
		fmt.Println(c.Request.Header.Get("Origin"))
		if strings.Index(c.Request.Header.Get("Origin"), setting.Domain) >= 0 {
			c.Writer.Header().Add("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
			passed = true
		}
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization,authorization,userid,username,x_orgId,x_orgname")
		c.Writer.Header().Add("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" && passed {
			c.JSON(200, "")
			return
		}
	}
}
