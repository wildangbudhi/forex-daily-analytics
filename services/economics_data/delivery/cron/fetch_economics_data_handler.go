package cron

import "log"

func (handler *EconomicsDataCronHandler) FetchEconomicsData() {

	log.Println("Fetch Economics Data From tradingeconomics.com")

	err := handler.economicsDataUsecase.FetchEconomicsData()

	if err != nil {
		log.Println("Error Fetch Economics Data From tradingeconomics.com : " + err.Error())
	}
}
