package main

import (
	"avatar-api-gin/controllers/avatarController"
	"avatar-api-gin/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	models.ConnectDatabase()
	models.ConnectRedis()

	r.GET("/api/avatar", avatarController.Index)
	r.GET("/api/avatar/:id", avatarController.Show)
	r.GET("/api/avatar/random", avatarController.Random)
	r.POST("/api/avatar", avatarController.Create)
	r.PUT("/api/avatar/:id", avatarController.Update)
	r.DELETE("/api/avatar/:id", avatarController.Delete)

	r.Run(":8000")
}
