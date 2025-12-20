package main

import (
	"auto-passport-api/controllers"
	"auto-passport-api/db"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func main() {

	db := db.New()
	router := gin.Default()
	app := controllers.PassportController{DB: db}

	router.POST("/api/v1/passports", app.PostPassports)

	if err := router.Run(":8080"); err != nil {
		slog.Error(fmt.Sprintf("failed to run server: %v", err))
	}

	slog.Info("Closing database connection")

	defer db.Pool.Close()
}
