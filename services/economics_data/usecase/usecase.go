package usecase

import (
	"github.com/wildangbudhi/forex-daily-analytics/services/economics_data/domain"
	"github.com/wildangbudhi/forex-daily-analytics/services/economics_data/domain/model"
)

type economicsDataUsecase struct {
	countryRepository                model.CountryRepository
	economicsDataCategoryRepository  model.EconomicsDataCategoryRepository
	economicsDataIndicatorRepository model.EconomicsDataIndicatorRepository
	economicsDataRepository          model.EconomicsDataRepository
}

func NewEconomicsDataUsecase(countryRepository model.CountryRepository, economicsDataCategoryRepository model.EconomicsDataCategoryRepository, economicsDataIndicatorRepository model.EconomicsDataIndicatorRepository, economicsDataRepository model.EconomicsDataRepository,
) domain.EconomicsDataUsecase {
	return &economicsDataUsecase{
		countryRepository:                countryRepository,
		economicsDataCategoryRepository:  economicsDataCategoryRepository,
		economicsDataIndicatorRepository: economicsDataIndicatorRepository,
		economicsDataRepository:          economicsDataRepository,
	}
}
