package pgsql

import (
	"github.com/go-pg/pg"
	"github.com/wildangbudhi/forex-daily-analytics/services/economics_data/domain/model"
)

type economicsDataIndicatorRepository struct {
	db *pg.DB
}

func NewEconomicDataIndicatorRepository(db *pg.DB) model.EconomicsDataIndicatorRepository {
	return &economicsDataIndicatorRepository{
		db: db,
	}
}

func (edir *economicsDataIndicatorRepository) GetByIndicatorName(indicatorName string) (*model.EconomicsDataIndicator, error) {

	economicsDataIndicator := &model.EconomicsDataIndicator{}

	err := edir.db.Model(economicsDataIndicator).Where("indicator_name = ?", indicatorName).Select()

	if err != nil {
		return nil, err
	}

	return economicsDataIndicator, nil

}

func (edir *economicsDataIndicatorRepository) Insert(indicatorName string) (int, error) {

	economicsDataIndicator := &model.EconomicsDataIndicator{
		IndicatorName: indicatorName,
	}

	res, err := edir.db.Model(economicsDataIndicator).Insert()

	return res.RowsAffected(), err

}
