package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (hs *HTTPServer) healthHandler(c *gin.Context) {
	if c.Request.Method != http.MethodGet || c.Request.URL.Path != "/api/health" {
		return
	}
	ok, result := hs.health.HealthCheck()
	c.Writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if ok {
		c.Writer.WriteHeader(200)
	} else {
		c.Writer.WriteHeader(500)
	}
	resultBytes, _ := json.Marshal(result)
	c.Writer.Write(resultBytes)
}
