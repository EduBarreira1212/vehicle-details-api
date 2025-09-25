package models

import (
	"errors"
	"regexp"
	"strings"
)

type FipeRequest struct {
	Plate string `json:"plate,omitempty"`
}

func (fipeRequest *FipeRequest) Validate() error {
	if err := fipeRequest.format(); err != nil {
		return err
	}

	oldFormat := regexp.MustCompile(`^[A-Z]{3}[0-9]{4}$`)

	mercosulFormat := regexp.MustCompile(`^[A-Z]{3}[0-9][A-Z0-9][0-9]{2}$`)

	if oldFormat.MatchString(fipeRequest.Plate) || mercosulFormat.MatchString(fipeRequest.Plate) {
		return nil
	}

	return errors.New("invalid plate format")
}

func (fipeRequest *FipeRequest) format() error {
	if fipeRequest.Plate == "" {
		return errors.New("plate is empty")
	}

	plate := strings.ToUpper(fipeRequest.Plate)
	plate = strings.ReplaceAll(plate, " ", "")
	plate = strings.ReplaceAll(plate, "-", "")

	fipeRequest.Plate = plate
	return nil
}
