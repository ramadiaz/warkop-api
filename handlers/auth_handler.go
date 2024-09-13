package handlers

import (
	"net/http"
	"warkop-api/dto"

	"github.com/gin-gonic/gin"
)

func (h *compHandlers) GenerateAPIKey(c *gin.Context) {
	name := c.Request.FormValue("name")
	secret := c.Request.FormValue("secret")

	key, err := h.service.GenerateAPIKey(name, secret)
	if err != nil {
		if err.Error() == "invalid secret" {
			c.JSON(http.StatusUnauthorized, dto.Response{Status: http.StatusUnauthorized, Error: err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "User successfully registerd", Body: key})
}
