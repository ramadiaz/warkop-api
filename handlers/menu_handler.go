package handlers

import (
	"net/http"
	"warkop-api/dto"

	"github.com/gin-gonic/gin"
)

func (h *compHandlers) RegisterMenu(c *gin.Context) {
	var data dto.Menu

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}

	err = h.service.RegisterMenu(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "Successfully register menu"})
}

func (h *compHandlers) GetAllMenu(c *gin.Context) {
	data, err := h.service.GetAllMenu()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "Data retrieved successfully", Body: data})
}
