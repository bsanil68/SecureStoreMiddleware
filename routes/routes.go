package routes

import (
	"SecureStore/controllers"
	"SecureStore/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, uploadController *controllers.UploadController) {
	router.Use(middleware.CORSMiddleware())
	router.POST("/Upload", uploadController.HandleUpload)
	router.GET("/ListFiles", uploadController.ListFiles)
	router.GET("/ViewFile", uploadController.GetFile)
}
