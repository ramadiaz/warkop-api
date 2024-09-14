package handlers

import (
	"database/sql"
	"net/http"
	"warkop-api/dto"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *compHandlers) RegisterUser(c *gin.Context) {
	var data dto.User

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}

	err = h.service.RegisterUser(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "Verification email has been sent"})
}

func (h *compHandlers) ResendVerification(c *gin.Context) {
	username := c.Query("un")

	err := h.service.GenerateVerificationEmail(username)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
            if pqErr.Code == "23505" {
                c.JSON(http.StatusConflict, dto.Response{Status: http.StatusConflict, Error: "Username already exist"})
            }
        }
		c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "New email verification sended"})
}

func (h *compHandlers) VerifyAccount(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, dto.Response{Status: http.StatusBadRequest, Error: "Token is required"})
		return
	}

	err := h.service.VerifyAccount(token)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, dto.Response{Status: http.StatusNotFound, Error: "Invalid token"})
			return
		} else if err.Error() == "410" {
			c.JSON(http.StatusGone, dto.Response{Status: http.StatusGone, Error: "Token expired"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "Email Verified"})
}

func (h *compHandlers) LoginUser(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	token, err := h.service.LoginUser(username, password)
	if err != nil {
		if err.Error() == "404" {
			c.JSON(http.StatusNotFound, dto.Response{Status: http.StatusNotFound, Error: "User not found"})
			return
		} else if err.Error() == "401" {
			c.JSON(http.StatusUnauthorized, dto.Response{Status: http.StatusUnauthorized, Error: "Invalid username or password"})
			return
		} else if err.Error() == "403" {
			c.JSON(http.StatusUnauthorized, dto.Response{Status: http.StatusUnauthorized, Error: "Your email is not verified"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
		}
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "Successfully login", Body: token})
}
