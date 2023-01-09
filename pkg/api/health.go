package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetServerHealthExample godoc
//	@Summary	get health example
//	@Schemes
//	@Description	do get health
//	@Tags			system
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Helloworld
//	@Router			/ [get]
func (hs *HTTPServer) healthHandler(c *gin.Context) {
	result := hs.health.HealthCheck()
	c.Writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	c.JSON(http.StatusOK, result)
}
