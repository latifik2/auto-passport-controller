package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/latifik2/auto-passport-controller/types"
)

type Response struct {
	Status  string
	Message string
}

func PostPassports(c *gin.Context) {
	var passportsBatch []types.CommonPassport

	if err := c.BindJSON(&passportsBatch); err != nil {
		slog.Error(fmt.Sprintf("Failed to unmarshal JSON body: %s", err))
		c.JSON(http.StatusBadRequest, Response{Status: "fail", Message: "Failed to unmarshal JSON body"})
		return
	}

	// implement logic to store data in PostgreSQL
	slog.Info(fmt.Sprintf("Successfully posted new passports: %s", passportsBatch))
	c.JSON(http.StatusCreated, Response{Status: "ok", Message: "Successfully posted new passports"})

}

// func insertPassportsBatch(passportsBatch *[]types.CommonPassport) {

// }
