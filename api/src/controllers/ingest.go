package controllers

import (
	"auto-passport-api/utils"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/latifik2/auto-passport-controller/types"
)

func (p *PassportController) PostPassports(c *gin.Context) {
	var passportsBatch []types.CommonPassport

	bodyBytes, err := c.GetRawData()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to GetRawData of JSON body: %s", err))
		c.JSON(http.StatusBadRequest, types.Response{Status: "fail", Message: "Failed to GetRawData of JSON body"})
		return
	}

	if len(bodyBytes) == 0 {
		slog.Error("Recieved empty body", "body", bodyBytes)
		c.JSON(http.StatusCreated, types.Response{Status: "fail", Message: "Recieved empty body"})
		return
	}

	// if err := c.BindJSON(&passportsBatch); err != nil {
	// 	slog.Error(fmt.Sprintf("Failed to unmarshal JSON body: %s", err))
	// 	c.JSON(http.StatusBadRequest, types.Response{Status: "fail", Message: "Failed to unmarshal JSON body"})
	// 	return
	// }

	if err := json.Unmarshal(bodyBytes, &passportsBatch); err != nil {
		slog.Error("Request contains invalid 'CommonPassport' JSON struct", "err", err)
		c.JSON(http.StatusBadRequest, types.Response{Status: "fail", Message: "Request contains invalid 'CommonPassport' JSON struct"})
		return
	}

	hashString := utils.GetSHA128String(bodyBytes)

	// implement logic to store data in PostgreSQL
	isDuplicatesHashes, err := p.DB.IsDuplicateHashes(hashString)

	if err != nil {
		slog.Error("Skipping insertion of raw JSON data due to searching for duplicates error")
		c.JSON(http.StatusInternalServerError, types.Response{Status: "fail", Message: "Skipping insertion of raw JSON data due to searching for duplicates error"})
		return
	}

	if isDuplicatesHashes {
		slog.Info("Found JSON with the same snapshop_hash in the database. Skipping the insertion into DB")
		c.JSON(http.StatusAlreadyReported, types.Response{Status: "ok", Message: "Found JSON with the same snapshop_hash in the database. Skipping the insertion into DB"})
		return
	} else {
		p.DB.InsertRawJSON(hashString, bodyBytes)

		slog.Info(fmt.Sprintf("Successfully posted new passports: %s", passportsBatch))
		c.JSON(http.StatusCreated, types.Response{Status: "ok", Message: "Successfully posted new passports"})
	}

}

// func insertPassportsBatch(passportsBatch *[]types.CommonPassport) {

// }
