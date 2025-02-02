package main

import (
	"github.com/DiegoDev2/img-processing/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run(":8080")
}
