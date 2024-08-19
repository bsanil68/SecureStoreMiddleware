// main.go
package main

import (
	"SecureStore/config"
	"SecureStore/controllers"
	"SecureStore/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	storjConfig := config.NewStorjConfig()
	uploadController := controllers.NewUploadController(storjConfig)

	routes.SetupRoutes(router, uploadController)

	router.Run(":8080")
}
