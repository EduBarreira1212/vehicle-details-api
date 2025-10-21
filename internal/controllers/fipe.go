package controllers

import (
	"net/http"
	"strconv"

	"github.com/EduBarreira1212/vehicle-details-api/internal/auth"
	"github.com/EduBarreira1212/vehicle-details-api/internal/models"
	"github.com/EduBarreira1212/vehicle-details-api/internal/responses"
	"github.com/EduBarreira1212/vehicle-details-api/internal/services"
	"github.com/gin-gonic/gin"
)

func GetFipe(c *gin.Context) {
	userIDInToken, _ := auth.GetUserIDFromContext(c)

	var plate models.FipeRequest
	if err := c.ShouldBindJSON(&plate); err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	if err := plate.Validate(); err != nil {
		responses.Error(c.Writer, http.StatusBadRequest, err)
		return
	}

	fipeDetails, err := services.GetFipe(c, userIDInToken, plate.Plate)
	if err != nil {
		responses.Error(c.Writer, http.StatusInternalServerError, err)
		return
	}

	outFipe := make([]models.Fipe, 0, len(fipeDetails.Fipe))
	for _, item := range fipeDetails.Fipe {
		outFipe = append(outFipe, models.Fipe{
			Brand:          item.Marca,
			Model:          item.Modelo,
			YearModel:      strconv.Itoa(item.AnoModelo),
			ReferenceMonth: item.MesReferencia,
			Fuel:           item.Combustivel,
			Value:          item.Valor,
		})
	}

	out := models.VehicleDataFiltered{
		Brand:     fipeDetails.InformacoesVeiculo.Marca,
		Model:     fipeDetails.InformacoesVeiculo.Modelo,
		Year:      fipeDetails.InformacoesVeiculo.Ano,
		ModelYear: fipeDetails.InformacoesVeiculo.AnoModelo,
		Color:     fipeDetails.InformacoesVeiculo.Cor,
		Chassis:   fipeDetails.InformacoesVeiculo.Chassi,
		City:      fipeDetails.InformacoesVeiculo.Municipio,
		State:     fipeDetails.InformacoesVeiculo.UF,
		Plate:     fipeDetails.InformacoesVeiculo.Placa,
		Fipe:      outFipe,
	}

	responses.JSON(c.Writer, http.StatusOK, out)
}
