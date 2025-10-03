package models

type VehicleDataFiltered struct {
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
