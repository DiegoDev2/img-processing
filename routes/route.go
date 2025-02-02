package routes

import (
	"net/http"

	"github.com/DiegoDev2/img-processing/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.Engine) {
	g.Use(middleware.RateLimiter)

	g.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "img-processing")
	})
}
