package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingExample godoc
//	@Summary	ping example
//	@Schemes
//	@Description	do ping
//	@Tags			system
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Helloworld
//	@Router			/ [get]
func (hs *HTTPServer) test(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
