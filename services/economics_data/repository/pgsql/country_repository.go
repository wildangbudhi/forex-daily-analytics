package pgsql

import (
	"github.com/go-pg/pg"
	"github.com/wildangbudhi/forex-daily-analytics/services/economics_data/domain/model"
)

type countryRepository struct {
	db *pg.DB
}

func NewCountryRepository(db *pg.DB) model.CountryRepository {
	return &countryRepository{
		db: db,
	}
}

func (cr *countryRepository) Fetch() ([]model.Country, error) {

	var countryList []model.Country

	err := cr.db.Model(&countryList).Select()

	if err != nil {
		return nil, err
	}

	return countryList, nil

}
