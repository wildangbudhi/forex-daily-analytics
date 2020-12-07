package pgsql

import (
	"github.com/go-pg/pg"
	"github.com/wildangbudhi/forex-daily-analytics/services/economics_data/domain/model"
)

type economicsDataCategoryRepository struct {
	db *pg.DB
}

func NewEconomicsDataCategoryRepository(db *pg.DB) model.EconomicsDataCategoryRepository {
	return &economicsDataCategoryRepository{
		db: db,
	}
}

func (edcr *economicsDataCategoryRepository) GetByCategoryName(categoryName string) (*model.EconomicsDataCategory, error) {

	economicsDataCategory := &model.EconomicsDataCategory{}

	err := edcr.db.Model(economicsDataCategory).Where("category_name = ?", categoryName).Select()

	if err != nil {
		return nil, err
	}

	return economicsDataCategory, nil
}

func (edcr *economicsDataCategoryRepository) Insert(categoryName string) (int, error) {

	economicsDataCategory := &model.EconomicsDataCategory{
		CategoryName: categoryName,
	}

	res, err := edcr.db.Model(economicsDataCategory).Insert()

	return res.RowsAffected(), err

}
