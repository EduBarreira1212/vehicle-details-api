package models

type Response struct {
	Codigo             int                `json:"codigo"`
	Msg                string             `json:"msg"`
	Fipe               []FipeItem         `json:"fipe"`
	InformacoesVeiculo InformacoesVeiculo `json:"informacoes_veiculo"`
	Tempo              int                `json:"tempo"`
	UndiadeTempo       string             `json:"undiade_tempo"`
	UnidadeTempo       string             `json:"unidade_tempo"`
	Algoritmo          string             `json:"algoritmo"`
	Placa              string             `json:"placa"`
}

type FipeItem struct {
	Similaridade     string `json:"similaridade"`
	Correspondencia  string `json:"correspondencia"`
	Marca            string `json:"marca"`
	Modelo           string `json:"modelo"`
	AnoModelo        int    `json:"ano_modelo"`
	CodigoFipe       string `json:"codigo_fipe"`
	CodigoMarca      string `json:"codigo_marca"`
	CodigoModelo     string `json:"codigo_modelo"`
	MesReferencia    string `json:"mes_referencia"`
	Combustivel      string `json:"combustivel"`
	Valor            string `json:"valor"`
	Desvalorizometro string `json:"desvalorizometro"`
	UnidadeValor     string `json:"unidade_valor"`
}

type InformacoesVeiculo struct {
	Marca       string `json:"marca"`
	Modelo      string `json:"modelo"`
	Ano         string `json:"ano"`
	AnoModelo   string `json:"ano_modelo"`
	Cor         string `json:"cor"`
	Chassi      string `json:"chassi"`
	Motor       string `json:"motor"`
	Municipio   string `json:"municipio"`
	UF          string `json:"uf"`
	Segmento    string `json:"segmento"`
	SubSegmento string `json:"sub_segmento"`
	Placa       string `json:"placa"`
}
