package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// ResponseModel representa o modelo da resposta retornada
type ResponseModel struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type CEPRepository struct {
	OrchestratorApiHost string
	OrchestratorApiPort string
}

func NewCEPRepository(orchestrator_api_host string, orchestrator_api_port string) *CEPRepository {
	return &CEPRepository{
		OrchestratorApiHost: orchestrator_api_host,
		OrchestratorApiPort: orchestrator_api_port,
	}
}

func (r *CEPRepository) IsValid(cep_address string) bool {
	check, _ := regexp.MatchString("^[0-9]{8}$", cep_address)
	return (len(cep_address) == 8 && cep_address != "" && check)
}

func (r *CEPRepository) Get(cep_address string) (*ResponseModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(
		"http://%s:%s/cep/%s",
		r.OrchestratorApiHost,
		r.OrchestratorApiPort,
		cep_address),
		nil,
	)
	if err != nil {
		log.Printf("Fail to create the request: %v", err)
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport,
			otelhttp.WithSpanNameFormatter(func(_ string, req *http.Request) string {
				return "get-cep-temp"
			}),
		),
	}

	resp, err := client.Do(req)

	log.Println("Requesting CEP", resp.Body)
	if err != nil {
		log.Printf("Fail to make the request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	ctx_err := ctx.Err()
	if ctx_err != nil {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			log.Printf("Max timeout reached: %v", err)
			return nil, err
		}
	}

	// Deserializar a resposta
	var responseModel ResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&responseModel); err != nil {
		log.Printf("Failed to decode response: %v", err)
		return nil, err
	}

	log.Printf("Response received: %+v", responseModel)

	// Retornar o modelo deserializado
	return &responseModel, nil

	// return nil
}
