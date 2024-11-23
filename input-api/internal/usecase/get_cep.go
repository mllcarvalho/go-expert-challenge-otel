package usecase

import (
	"github.com/mllcarvalho/go-expert-challenge-otel/input-api/internal/entity"
)

type CEPInputDTO struct {
	CEP string `json:"cep"`
}

type CEPOutputDTO struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
}

type WeatherOutputDTO struct {
	City       string  `json:"city"`
	Celcius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type GetCEPUseCase struct {
	CEPRepository entity.CEPRepositoryInterface
}

func NewGetCEPUseCase(cep_repository entity.CEPRepositoryInterface) *GetCEPUseCase {
	return &GetCEPUseCase{
		CEPRepository: cep_repository,
	}
}

func (c *GetCEPUseCase) Execute(input CEPInputDTO) (*WeatherOutputDTO, error) {
	// Chamar o método Get e capturar os dois valores retornados
	responseModel, err := c.CEPRepository.Get(input.CEP)
	if err != nil {
		return nil, err
	}

	// Mapear os dados retornados para o DTO de saída
	output := &WeatherOutputDTO{
		City:       responseModel.City,
		Celcius:    responseModel.TempC,
		Fahrenheit: responseModel.TempF,
		Kelvin:     responseModel.TempK,
	}

	return output, nil
}
