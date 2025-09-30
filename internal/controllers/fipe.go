package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/EduBarreira1212/vehicle-details-api/internal/models"
	"github.com/EduBarreira1212/vehicle-details-api/internal/responses"
	"github.com/EduBarreira1212/vehicle-details-api/internal/services"
	"github.com/gin-gonic/gin"
)

func GetFipe(c *gin.Context) {
	parameters := c.Param("userID")

	ID, err := strconv.ParseUint(parameters, 10, 64)
	if err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		responses.Error(c.Writer, http.StatusUnprocessableEntity, err)
		return
	}

	var plate models.FipeRequest
	if err = json.Unmarshal(body, &plate); err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	if err = plate.Validate(); err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	fipeDetails, err := services.GetFipe(c, ID, plate.Plate)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(c.Writer, http.StatusCreated, fipeDetails)
}
