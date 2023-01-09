package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/magicLian/gostarter/pkg/models"
)

type AuthOptions struct {
	RequireSignedIn bool
	RequireAdmin    bool
}

func accessForbidden(c *models.ReqContext) {
	c.Error(403, "Permission denied", nil)
	return
}

func notAuthorized(c *models.ReqContext) {
	c.Error(401, "Unauthorized", nil)
	return
}

func Auth(options *AuthOptions) gin.HandlerFunc {
	return func(context *gin.Context) {
		c := models.CtxGetByGin(context)
		if !c.IsSignedIn {
			notAuthorized(c)
			return
		}
		if options.RequireAdmin {
			if !c.IsAdmin {
				accessForbidden(c)
				return
			}
		}
	}
}
