package cron

import (
	"github.com/robfig/cron/v3"
	"github.com/wildangbudhi/forex-daily-analytics/services/economics_data/domain"
)

type EconomicsDataCronHandler struct {
	economicsDataUsecase domain.EconomicsDataUsecase
}

func NewEconomicsDataCronHandler(scheduler *cron.Cron, economicsDataUsecase domain.EconomicsDataUsecase) {

	handler := EconomicsDataCronHandler{
		economicsDataUsecase: economicsDataUsecase,
	}

	scheduler.AddFunc("45 * * * *", handler.FetchEconomicsData)

}
