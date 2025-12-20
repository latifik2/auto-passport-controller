package controllers

import (
	"auto-passport-api/utils"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/latifik2/auto-passport-controller/types"
)

func (p *PassportController) PostPassports(c *gin.Context) {
	var passportsBatch []types.CommonPassport

	if err := c.BindJSON(&passportsBatch); err != nil {
		slog.Error(fmt.Sprintf("Failed to unmarshal JSON body: %s", err))
		c.JSON(http.StatusBadRequest, Response{Status: "fail", Message: "Failed to unmarshal JSON body"})
		return
	}

	hashBytes, err := c.GetRawData()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to GetRawData of JSON body: %s", err))
		c.JSON(http.StatusBadRequest, Response{Status: "fail", Message: "Failed to GetRawData of JSON body"})
		return
	}

	fmt.Println(hashBytes)

	hashString := utils.GetSHA128String(hashBytes)

	// implement logic to store data in PostgreSQL
	p.DB.InsertRawJSON(hashString, hashBytes)

	slog.Info(fmt.Sprintf("Successfully posted new passports: %s", passportsBatch))
	c.JSON(http.StatusCreated, Response{Status: "ok", Message: "Successfully posted new passports"})

}

// func insertPassportsBatch(passportsBatch *[]types.CommonPassport) {

// }
