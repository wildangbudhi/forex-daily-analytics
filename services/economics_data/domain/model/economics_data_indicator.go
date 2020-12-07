package model

type EconomicsDataIndicator struct {
	Id            int    `json:"id"`
	IndicatorName string `json:"indicator_name"`
}

type EconomicsDataIndicatorRepository interface {
	GetByIndicatorName(indicatorName string) (*EconomicsDataIndicator, error)
	Insert(indicatorName string) (int, error)
}
