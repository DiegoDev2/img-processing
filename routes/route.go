package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DiegoDev2/img-processing/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.Engine) {

	g.Use(middleware.RateLimiter)

	g.MaxMultipartMemory = 50 << 20 // 50 MB
	g.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("image")
		log.Println(file.Filename)
		log.Println(file.Size)
		log.Println(file.Header)

		c.SaveUploadedFile(file, file.Filename)
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
}
