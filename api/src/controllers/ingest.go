package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func PostPassports(c *gin.Context) {
	resp := Response{
		Status:  "ok",
		Message: "succesful ingest",
	}

	c.IndentedJSON(http.StatusOK, resp)
}
