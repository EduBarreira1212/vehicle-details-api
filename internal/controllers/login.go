package controllers

import (
	"net/http"

	"github.com/EduBarreira1212/vehicle-details-api/internal/auth"
	"github.com/EduBarreira1212/vehicle-details-api/internal/config"
	"github.com/EduBarreira1212/vehicle-details-api/internal/models"
	"github.com/EduBarreira1212/vehicle-details-api/internal/repositories"
	"github.com/EduBarreira1212/vehicle-details-api/internal/responses"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	repository := repositories.NewUserRepository(config.DB)
	userSaved, err := repository.GetByEmail(c.Request.Context(), user.Email)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	if err = auth.VerifyPassword(userSaved.Password, user.Password); err != nil {
		responses.Error(c.Writer, http.StatusUnauthorized, err)
		return
	}

	token, err := auth.CreateToken(userSaved.ID)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(c.Writer, http.StatusCreated, token)
}
