package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/EduBarreira1212/vehicle-details-api/internal/config"
	"github.com/EduBarreira1212/vehicle-details-api/internal/models"
	"github.com/EduBarreira1212/vehicle-details-api/internal/repositories"
	"github.com/gin-gonic/gin"
)

type GetFipePriceRequest struct {
	Placa string `json:"placa"`
	Token string `json:"token"`
}

func GetFipe(c *gin.Context, userID uint64, plate string) (models.Response, error) {
	url := "https://api.placafipe.com.br/getplacafipe"

	payload := GetFipePriceRequest{
		Placa: plate,
		Token: config.LoadConfig().FIPE_API_TOKEN,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return models.Response{}, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return models.Response{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Response{}, err
	}

	repository := repositories.NewUserRepository(config.DB)
	err = repository.UpdateUserHistory(c.Request.Context(), userID, plate)
	if err != nil {
		return models.Response{}, err
	}

	var respJSON models.Response
	if err := json.Unmarshal(body, &respJSON); err != nil {
		return models.Response{}, err
	}

	return respJSON, nil
}
