package entity

import "github.com/mllcarvalho/go-expert-challenge-otel/input-api/internal/infra/repo"

type CEPRepositoryInterface interface {
	Get(string) (*repo.ResponseModel, error)
	IsValid(string) bool
}
