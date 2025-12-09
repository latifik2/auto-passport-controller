package main

import (
	"auto-passport-api/controllers"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/api/v1/passports", controllers.PostPassports)

	if err := router.Run("localhost:8080"); err != nil {
		slog.Error(fmt.Sprintf("failed to run server: %v", err))
	}
}
