package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/latifik2/auto-passport-controller/types"
)

func (p *PassportController) GetPassports(c *gin.Context) {
	var passportsBatch []types.CommonPassport

	passportsBytes, err := p.DB.SelectActualPassports()

	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{Status: "fail", Message: err.Error()})
		return
	}

	if err := json.Unmarshal(passportsBytes, &passportsBatch); err != nil {
		slog.Error("Failed to unmarshal data 'from passports_raw' table into '[]types.CommonPassport'", "err", err)
		c.JSON(http.StatusInternalServerError, types.Response{Status: "fail", Message: "Internal Server Error occured. Please, contact administrator"})
		return
	}

	c.JSON(http.StatusOK, passportsBatch)
}
