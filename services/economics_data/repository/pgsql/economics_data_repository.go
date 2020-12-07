package pgsql

import (
	"database/sql"
	"time"

	"github.com/go-pg/pg"
	"github.com/wildangbudhi/forex-daily-analytics/services/economics_data/domain/model"
)

type economicsDataRepository struct {
	db *pg.DB
}

func NewEconomicsDataRepository(db *pg.DB) model.EconomicsDataRepository {
	return &economicsDataRepository{
		db: db,
	}
}

func (edr *economicsDataRepository) Insert(datetime time.Time, countryID string, economicsDataCategoryID int, economicsDataIndicatorID int, lastValue float64, previousValue float64) (int, error) {

	economicsData := &model.EconomicsData{
		Datetime:                 datetime,
		CountryID:                countryID,
		EconomicsDataCategoryID:  economicsDataCategoryID,
		EconomicsDataIndicatorID: economicsDataIndicatorID,
		LastValue:                sql.NullFloat64{Float64: lastValue, Valid: true},
		PreviousValue:            sql.NullFloat64{Float64: previousValue, Valid: true},
	}

	res, err := edr.db.Model(economicsData).Insert()

	return res.RowsAffected(), err
}

func (edr *economicsDataRepository) GetEconomicsDataScore() ([]model.EconomicsDataScore, error) {

	queryString := `
	with econimic_data_with_score as (
		select
			c.currency_code,
			edc.category_name,
			edi.indicator_name,
			ed.last_value,
			ed.previous_value,
			case 
				when ed.economics_data_indicator_id IN (97, 17) then
					case
						when ed.last_value > ed.previous_value then -1
						when ed.last_value = ed.previous_value then 0
						when ed.last_value < ed.previous_value then 1
					end
				else
					case
						when ed.last_value > ed.previous_value then 1
						when ed.last_value = ed.previous_value then 0
						when ed.last_value < ed.previous_value then -1
					end 
			end as score
		from 
			economics_data ed 
		LEFT JOIN economics_data_categories edc on ed.economics_data_category_id = edc.id 
		LEFT JOIN countries c on c.id = ed.country_id
		LEFT JOIN economics_data_indicators edi on ed.economics_data_indicator_id = edi .id 
		WHERE 
			ed.datetime = (SELECT DISTINCT datetime FROM economics_data ORDER BY datetime DESC LIMIT 1) and
			(
				ed.country_id = 'united-states' and 
				ed.economics_data_indicator_id in (1,2,17,18,50,51,79,80,81,97,98,114,115,116,117,175,176,192,108,109)
			) or
			(
				ed.country_id = 'euro-area' and 
				ed.economics_data_indicator_id in (1,2,17,50,51,79,80,81,97,98,114,115,117,175,176,108,109)
			) or
			(
				ed.country_id = 'japan' and 
				ed.economics_data_indicator_id in (1,2,241,17,50,51,79,80,81,97,98,114,115,117,175,176,108,109)
			) or
			(
				ed.country_id = 'united-kingdom' and 
				ed.economics_data_indicator_id in (1,2,17,50,51,79,80,81,97,98,114,115,117,175,176,108,109)
			) or
			(
				ed.country_id = 'canada' and 
				ed.economics_data_indicator_id in (1,2,241,17,50,51,79,80,81,97,98,114,115,175,176,108,109,192)
			) or
			(
				ed.country_id = 'australia' and 
				ed.economics_data_indicator_id in (1,2,17,50,51,79,80,81,97,98,114,115,117,175,176,108,109,192)
			) or
			(
				ed.country_id = 'switzerland' and 
				ed.economics_data_indicator_id in (1,2,17,50,51,79,80,81,97,98,114,115,175,176,108,109)
			) or
			(
				ed.country_id = 'new-zealand' and 
				ed.economics_data_indicator_id in (1,2,17,50,51,79,80,81,97,98,114,115,175,176,108,109,192)
			)
	),
	
	economic_data_with_score_per_category as (
		select 
			edws.currency_code,
			edws.category_name,
			case
				when sum( edws.score ) > 0 then 1
				when sum( edws.score ) = 0 then 0
				when sum( edws.score ) < 0 then -1
			end as score
		from econimic_data_with_score edws
		group by edws.category_name, edws.currency_code
		order by edws.currency_code
	)
	
	select 
		edwspc.currency_code,
		sum( edwspc.score ) as score
	from 
		economic_data_with_score_per_category edwspc
	group by
		edwspc.currency_code
	
	`

	var economicsDataScore []model.EconomicsDataScore

	_, err := edr.db.Query(&economicsDataScore, queryString)

	if err != nil {
		return nil, err
	}

	return economicsDataScore, err

}
