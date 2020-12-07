package depedencyinjection

import (
	"github.com/wildangbudhi/forex-daily-analytics/services/economics_data/delivery/cron"
	"github.com/wildangbudhi/forex-daily-analytics/services/economics_data/repository/pgsql"
	"github.com/wildangbudhi/forex-daily-analytics/services/economics_data/usecase"
	"github.com/wildangbudhi/forex-daily-analytics/utils"
)

func EconomicsData(server *utils.Server) {
	countryRepository := pgsql.NewCountryRepository(server.DB)
	economicsDataCategoryRepository := pgsql.NewEconomicsDataCategoryRepository(server.DB)
	economicsDataIndicatorRepository := pgsql.NewEconomicDataIndicatorRepository(server.DB)
	economicsDataRepository := pgsql.NewEconomicsDataRepository(server.DB)
	economicsDataUsecase := usecase.NewEconomicsDataUsecase(
		countryRepository,
		economicsDataCategoryRepository,
		economicsDataIndicatorRepository,
		economicsDataRepository,
	)
	cron.NewEconomicsDataCronHandler(server.Scheduler, economicsDataUsecase)
}
