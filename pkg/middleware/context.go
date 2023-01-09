package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/magicLian/gostarter/pkg/log"
	"github.com/magicLian/gostarter/pkg/models"
)

func (m *MiddleWare) InitContextHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCtx := &models.ReqContext{
			Log:        log.New("middleware"),
			IsSignedIn: false,
			Context:    c,
		}

		reqCtx.Context = c
		c.Set("ReqContext", reqCtx)
	}
}
