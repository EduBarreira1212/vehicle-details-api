package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/EduBarreira1212/vehicle-details-api/internal/auth"
	"github.com/EduBarreira1212/vehicle-details-api/internal/config"
	"github.com/EduBarreira1212/vehicle-details-api/internal/models"
	"github.com/EduBarreira1212/vehicle-details-api/internal/repositories"
	"github.com/EduBarreira1212/vehicle-details-api/internal/responses"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	if err := user.Prepare("register"); err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	repository := repositories.NewUserRepository(config.DB)
	userCreated, err := repository.Create(c.Request.Context(), user.Name, user.Email, user.Password)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	token, err := auth.CreateToken(userCreated.ID)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(c.Writer, http.StatusCreated, token)
}

func GetMyProfile(c *gin.Context) {
	userIDInToken, _ := auth.GetUserIDFromContext(c)

	repository := repositories.NewUserRepository(config.DB)
	user, err := repository.GetById(c.Request.Context(), userIDInToken)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(c.Writer, http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	parameters := c.Param("userID")

	ID, err := strconv.ParseUint(parameters, 10, 64)
	if err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	repository := repositories.NewUserRepository(config.DB)
	user, err := repository.GetById(c.Request.Context(), ID)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(c.Writer, http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	parameters := c.Param("userID")

	ID, err := strconv.ParseUint(parameters, 10, 64)
	if err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	userIDInToken, _ := auth.GetUserIDFromContext(c)

	if userIDInToken != ID {
		responses.Error(c.Writer, http.StatusForbidden, errors.New("isn't possible to update a different user"))
		return
	}

	var user models.User
	if err = c.ShouldBindJSON(&user); err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	if err = user.Prepare("update"); err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	repository := repositories.NewUserRepository(config.DB)
	err = repository.Update(c.Request.Context(), ID, user.Name, user.Email)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(c.Writer, http.StatusNoContent, nil)
}

func UpdatePassword(c *gin.Context) {
	parameters := c.Param("userID")

	ID, err := strconv.ParseUint(parameters, 10, 64)
	if err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	userIDInToken, _ := auth.GetUserIDFromContext(c)

	if userIDInToken != ID {
		responses.Error(c.Writer, http.StatusForbidden, errors.New("isn't possible to update a different user's password"))
		return
	}

	var password models.Password
	if err = c.ShouldBindJSON(&password); err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	repository := repositories.NewUserRepository(config.DB)
	passwordInDB, err := repository.GetPassword(c.Request.Context(), ID)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	if err = auth.VerifyPassword(passwordInDB, password.Current); err != nil {
		responses.Error(c.Writer, http.StatusUnauthorized, err)
		return
	}

	passwordWithHash, err := auth.Hash(password.New)
	if err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	err = repository.UpdatePassword(c.Request.Context(), ID, string(passwordWithHash))
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(c.Writer, http.StatusNoContent, nil)
}

func DeleteUser(c *gin.Context) {
	parameters := c.Param("userID")

	ID, err := strconv.ParseUint(parameters, 10, 64)
	if err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	userIDInToken, _ := auth.GetUserIDFromContext(c)

	if userIDInToken != ID {
		responses.Error(c.Writer, http.StatusForbidden, errors.New("isn't possible to delete a different user"))
		return
	}

	repository := repositories.NewUserRepository(config.DB)
	err = repository.Delete(c.Request.Context(), ID)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(c.Writer, http.StatusNoContent, nil)
}

func GetUserHistory(c *gin.Context) {
	parameters := c.Param("userID")

	ID, err := strconv.ParseUint(parameters, 10, 64)
	if err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	userIDInToken, _ := auth.GetUserIDFromContext(c)

	if userIDInToken != ID {
		responses.Error(c.Writer, http.StatusForbidden, errors.New("isn't possible to get a different user's history"))
		return
	}

	repository := repositories.NewUserRepository(config.DB)
	userHistory, err := repository.GetUserHistoryById(c.Request.Context(), ID)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(c.Writer, http.StatusOK, userHistory)
}
