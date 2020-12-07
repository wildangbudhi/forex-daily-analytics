package usecase

import "github.com/wildangbudhi/forex-daily-analytics/services/economics_data/domain/model"

func (edu *economicsDataUsecase) GetEconomicsDataScore() ([]model.EconomicsDataScore, error) {

	economicsDataScore, err := edu.economicsDataRepository.GetEconomicsDataScore()

	return economicsDataScore, err

}
