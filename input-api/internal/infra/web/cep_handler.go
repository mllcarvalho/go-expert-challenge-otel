package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mllcarvalho/go-expert-challenge-otel/configs"
	"github.com/mllcarvalho/go-expert-challenge-otel/input-api/internal/entity"
	"github.com/mllcarvalho/go-expert-challenge-otel/input-api/internal/infra/repo"
	"github.com/mllcarvalho/go-expert-challenge-otel/input-api/internal/usecase"
	"go.opentelemetry.io/otel/trace"
)

type WebCEPHandler struct {
	CEPRepository entity.CEPRepositoryInterface
	Configs       *configs.Conf
	Tracer        trace.Tracer
}

func NewWebCEPHandler(conf *configs.Conf, tracer trace.Tracer) *WebCEPHandler {
	return &WebCEPHandler{
		CEPRepository: repo.NewCEPRepository(conf.OrchestratorApiHost, conf.OrchestratorApiPort),
		Configs:       conf,
		Tracer:        tracer,
	}
}

func (h *WebCEPHandler) Get(w http.ResponseWriter, r *http.Request) {
	resp, err := io.ReadAll(r.Body)
	log.Println("resp: ", resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("fail to read the response: %v", err), http.StatusInternalServerError)
		return
	}

	var cep_data usecase.CEPInputDTO
	err = json.Unmarshal(resp, &cep_data)
	if err != nil {
		http.Error(w, fmt.Sprintf("fail to parse the cep_data: %v", err), http.StatusInternalServerError)
		return
	}

	validate_cep_dto := usecase.ValidateCEPInputDTO{
		CEP: cep_data.CEP,
	}

	validateCEP := usecase.NewValidateCEPUseCase(h.CEPRepository)
	is_valid := validateCEP.Execute(validate_cep_dto)
	if !is_valid {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	get_cep_dto := usecase.CEPInputDTO{
		CEP: cep_data.CEP,
	}

	getCEP := usecase.NewGetCEPUseCase(h.CEPRepository)
	cepOutput, err := getCEP.Execute(get_cep_dto)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting cep: %v", err), http.StatusInternalServerError)
		return
	}

	// // Use a nova assinatura de Get
	// responseModel, err := h.CEPRepository.Get(cep_data.CEP)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("error getting cep: %v", err), http.StatusInternalServerError)
	// 	return
	// }

	// Retornar o resultado como JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(cepOutput); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
