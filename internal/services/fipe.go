package services

import (
	"github.com/EduBarreira1212/vehicle-details-api/internal/config"
	"github.com/EduBarreira1212/vehicle-details-api/internal/repositories"
	"github.com/gin-gonic/gin"
)

func GetFipe(c *gin.Context, userID uint64, plate string) (string, error) {
	//call api

	repository := repositories.NewUserRepository(config.DB)
	err := repository.UpdateUserHistory(c.Request.Context(), userID, plate)
	if err != nil {
		return "", err
	}

	return "Ok!", nil
}
