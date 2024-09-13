package handlers

import (
	"net/http"
	"warkop-api/dto"
	"warkop-api/helpers"

	"github.com/gin-gonic/gin"
)

func (h *compHandlers) RegisterUser(c *gin.Context) {
	var data dto.User

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}

	token, err := h.service.RegisterUser(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "User successfully registerd", Body: token})
}

func (h *compHandlers) ResendVerification(c *gin.Context) {
	err := h.service.GenerateVerificationEmail(helpers.GetUserData(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "New email verification sended"})
}
