package domain

import "github.com/wildangbudhi/forex-daily-analytics/services/economics_data/domain/model"

type EconomicsDataUsecase interface {
	FetchEconomicsData() error
	GetEconomicsDataScore() ([]model.EconomicsDataScore, error)
}

type Rules struct {
	Category      string
	CategoryIndex int
	Function      func(float64, float64) int
}
