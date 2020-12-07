package model

import (
	"database/sql"
	"time"
)

type EconomicsData struct {
	Datetime                 time.Time               `json:"datetime"`
	CountryID                string                  `json:"country_id"`
	Country                  *Country                `json:"-" pg:"rel:has-one"`
	EconomicsDataCategoryID  int                     `json:"economics_data_category_id"`
	EconomicsDataCategory    *EconomicsDataCategory  `json:"-" pg:"rel:has-one"`
	EconomicsDataIndicatorID int                     `json:"economics_data_indicator_id"`
	EconomicsDataIndicator   *EconomicsDataIndicator `json:"-" pg:"rel:has-one"`
	LastValue                sql.NullFloat64         `json:"last_value" pg:",use_zero"`
	PreviousValue            sql.NullFloat64         `json:"previous_value" pg:",use_zero"`
}

type EconomicsDataScore struct {
	CurrencyCode string  `json:"currency_code"`
	Score        float64 `json:"score"`
}

type EconomicsDataRepository interface {
	Insert(datetime time.Time, countryID string, economicsDataCategoryID int, economicsDataIndicatorID int, lastValue float64, previousValue float64) (int, error)
	GetEconomicsDataScore() ([]EconomicsDataScore, error)
}
