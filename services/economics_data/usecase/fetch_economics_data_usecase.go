package usecase

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type tradingEconomicsData struct {
	Category  string
	Indicator string
	Last      float64
	Previous  float64
}

func fetchData(countryCode string) []tradingEconomicsData {

	res, err := http.Get("https://tradingeconomics.com/" + countryCode + "/indicators")

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	tradingEconomicData := make([]tradingEconomicsData, 0)
	tabs := []string{"gdp", "labour", "prices", "money", "trade", "government", "business", "consumer", "health"}

	for _, tabName := range tabs {

		doc.Find("#" + tabName).Children().Each(func(i int, tabContent *goquery.Selection) {
			if i == 1 || i == 3 {

				var category string

				tabContent.Children().Children().Each(func(j int, tableContent *goquery.Selection) {

					if j == 0 {
						tableContent.Children().Children().Each(func(k int, thContent *goquery.Selection) {
							if k == 0 {
								category = strings.TrimSpace(thContent.Text())
							}
						})
					} else if j == 1 {
						tableContent.Children().Each(func(k int, trContent *goquery.Selection) {

							row := tradingEconomicsData{Category: category}

							trContent.Children().Each(func(l int, komp *goquery.Selection) {

								strKomp := strings.Replace(komp.Text(), "\n", "", -1)
								strKomp = regexp.MustCompile(`(\(.*\))`).ReplaceAllLiteralString(strKomp, "")
								strKomp = strings.TrimSpace(strKomp)

								if l == 0 {
									row.Indicator = strKomp
								} else if l == 1 {
									komponen, err := strconv.ParseFloat(strKomp, 64)

									if err != nil {
										row.Last = 0.0
									} else {
										row.Last = float64(komponen)
									}
								} else if l == 3 {
									komponen, err := strconv.ParseFloat(strKomp, 64)

									if err != nil {
										row.Previous = 0.0
									} else {
										row.Previous = float64(komponen)
									}
								}
							})

							tradingEconomicData = append(tradingEconomicData, row)

						})
					}
				})

			}
		})

	}

	return tradingEconomicData
}

func (edu *economicsDataUsecase) FetchEconomicsData() error {

	countries, err := edu.countryRepository.Fetch()

	if err != nil {
		return err
	}

	now := time.Now()

	for i := 0; i < len(countries); i++ {
		data := fetchData(countries[i].Id)

		for j := 0; j < len(data); j++ {

			economicsDataCategory, err := edu.economicsDataCategoryRepository.GetByCategoryName(data[j].Category)

			if err != nil {
				return err
			}

			economicsDataIndicator, err := edu.economicsDataIndicatorRepository.GetByIndicatorName(data[j].Indicator)

			if err != nil {
				return err
			}

			_, err = edu.economicsDataRepository.Insert(now, countries[i].Id, economicsDataCategory.Id, economicsDataIndicator.Id, data[j].Last, data[j].Previous)

			if err != nil {
				return err
			}

		}

	}

	return nil

}
