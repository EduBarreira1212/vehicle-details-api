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

type FipeInfoFiltered struct {
	Brand     string `json:"brand"`
	Model     string `json:"model"`
	Year      string `json:"year"`
	ModelYear string `json:"model_year"`
	Color     string `json:"color"`
	Chassis   string `json:"chassis"`
	City      string `json:"city"`
	State     string `json:"state"`
	Plate     string `json:"plate"`
	Fipe      []Fipe `json:"fipe"`
}

type Fipe struct {
	Brand          string `json:"brand"`
	Model          string `json:"model"`
	YearModel      string `json:"year_model"`
	ReferenceMonth string `json:"reference_month"`
	Fuel           string `json:"fuel"`
	Value          string `json:"value"`
}

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

	outFipe := make([]Fipe, 0, len(fipeDetails.Fipe))
	for _, item := range fipeDetails.Fipe {
		outFipe = append(outFipe, Fipe{
			Brand:          item.Marca,
			Model:          item.Modelo,
			YearModel:      strconv.Itoa(item.AnoModelo),
			ReferenceMonth: item.MesReferencia,
			Fuel:           item.Combustivel,
			Value:          item.Valor,
		})
	}

	out := FipeInfoFiltered{
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

	responses.JSON(c.Writer, http.StatusCreated, out)
}
