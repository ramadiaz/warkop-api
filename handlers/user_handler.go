package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"warkop-api/dto"

	"github.com/dgrijalva/jwt-go"
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
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				c.JSON(http.StatusConflict, dto.Response{Status: http.StatusConflict, Error: "Username already exist"})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "Verification email has been sent"})
}

func (h *compHandlers) ResendVerification(c *gin.Context) {
	username := c.Query("un")

	err := h.service.GenerateVerificationEmail(username)
	if err != nil {
		if err.Error() == "account already verified" {
			c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "Account already verified"})
			return
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
			c.JSON(http.StatusForbidden, dto.Response{Status: http.StatusForbidden, Error: "Your email is not verified"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "Successfully login", Body: token})
}

func (h *compHandlers) RequestResetPassword(c *gin.Context) {
	username := c.Query("un")

	data, err := h.service.RequestResetPassword(username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, dto.Response{Status: http.StatusNotFound, Error: "Username not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "OTP Code successfully sent to your email!", Body: data.ID})
}

func (h *compHandlers) VerifyResetPassword(c *gin.Context) {
	var data dto.OTPVerifyToken

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Status: http.StatusBadRequest, Error: err.Error()})
		return
	}

	result, err := h.service.VerifyResetPassword(data)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, dto.Response{Status: http.StatusNotFound, Error: "Token Invalid"})
			return
		} else if err.Error() == "410" {
			c.JSON(http.StatusGone, dto.Response{Status: http.StatusGone, Error: "Expired OTP"})
			return
		} else if err.Error() == "401" {
			c.JSON(http.StatusUnauthorized, dto.Response{Status: http.StatusUnauthorized, Error: "Invalid OTP"})
			return
		} else if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "22P02" {
				c.JSON(http.StatusNotFound, dto.Response{Status: http.StatusNotFound, Error: "Token Invalid"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "OTP Code verified", Body: result})
}

func (h *compHandlers) ResetPassword(c *gin.Context) {
	user_token := c.Query("token")
	password := c.Request.FormValue("password")

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error getting secret"})
		return
	}

	var secretKey = []byte(secret)

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(user_token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
		return
	}

	if !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
		return
	}

	user_data := dto.User{
		ID:       claims["id"].(string),
		Password: password,
	}

	err = h.service.ResetPassword(user_data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Status: http.StatusInternalServerError, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Status: http.StatusOK, Message: "Password successfully reseted"})

}
