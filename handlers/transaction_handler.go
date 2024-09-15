package handlers

import (
	"net/http"
	"warkop-api/dto"
	"warkop-api/helpers"

	"github.com/gin-gonic/gin"
)

func (h *compHandlers) RegisterTransaction(c *gin.Context) {
	var data dto.Transaction

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}

	user_data := helpers.GetUserData(c)

	data.CashierID = user_data.ID

	err = h.service.RegisterTransaction(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "Transaction recorded"})
}
