package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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

func greaterBetter(last, prev float64) int {
	if last > prev {
		return 1
	} else if last == prev {
		return 0
	} else {
		return -1
	}
}

func lessBetter(last, prev float64) int {
	if last < prev {
		return 1
	} else if last == prev {
		return 0
	} else {
		return -1
	}
}

var mapScore map[string]map[string]func(last, prev float64) int

func initMap() {

	mapScore = make(map[string]map[string]func(last float64, prev float64) int)

	mapScore["GDP"] = map[string]func(last float64, prev float64) int{
		"GDP Growth Rate":        greaterBetter,
		"GDP Annual Growth Rate": greaterBetter,
		"GDP Growth Annualized":  greaterBetter,
	}

	mapScore["Labour"] = map[string]func(last float64, prev float64) int{
		"Unemployment Rate": lessBetter,
		"Non Farm Payrolls": greaterBetter,
	}

	mapScore["Prices"] = map[string]func(last float64, prev float64) int{
		"Inflation Rate":     greaterBetter,
		"Inflation Rate Mom": greaterBetter,
	}

	mapScore["Trade"] = map[string]func(last float64, prev float64) int{
		"Balance of Trade":       greaterBetter,
		"Current Account":        greaterBetter,
		"Current Account to GDP": greaterBetter,
	}

	mapScore["Government"] = map[string]func(last float64, prev float64) int{
		"Government Debt to GDP": lessBetter,
		"Government Budget":      greaterBetter,
	}

	mapScore["Business"] = map[string]func(last float64, prev float64) int{
		"Business Confidence":   greaterBetter,
		"Manufacturing PMI":     greaterBetter,
		"Non Manufacturing PMI": greaterBetter,
		"Services PMI":          greaterBetter,
	}

	mapScore["Consumer"] = map[string]func(last float64, prev float64) int{
		"Consumer Confidence": greaterBetter,
		"Retail Sales MoM":    greaterBetter,
	}

	mapScore["Housing"] = map[string]func(last float64, prev float64) int{
		"Building Permits": greaterBetter,
	}

	mapScore["Taxes"] = map[string]func(last float64, prev float64) int{
		"Corporate Tax Rate":       greaterBetter,
		"Personal Income Tax Rate": greaterBetter,
	}

}

func calculateScore(economicsData []tradingEconomicsData) (int, float64) {

	var score map[string]map[string]int = make(map[string]map[string]int)
	var moneySupply float64

	for i := 0; i < len(economicsData); i++ {
		var category, indicator string = economicsData[i].Category, economicsData[i].Indicator
		var last, prev float64 = economicsData[i].Last, economicsData[i].Previous

		if category == "Money" && indicator == "Money Supply M1" {
			moneySupply = (last - prev) / last
		} else {

			indicatorValue, isIndicatorValueExist := mapScore[category]

			if isIndicatorValueExist {
				isBetterFunction, ok := indicatorValue[indicator]

				if ok {

					if len(score[category]) == 0 {
						score[category] = make(map[string]int)
					}

					score[category][indicator] = isBetterFunction(last, prev)
					// fmt.Printf("%s - %s = %f - %f --> %d\n", category, indicator, last, prev, score[category][indicator])
				}
			}
		}

	}

	var scoreByCategoryList []int

	for _, categoryValue := range score {

		totalScore := 0

		for _, indicatorValue := range categoryValue {
			totalScore += indicatorValue
		}

		// fmt.Printf("%s %d\n", categoryName, totalScore)

		scoreByCategoryList = append(scoreByCategoryList, totalScore)

	}

	var finalScore int = 0

	for i := 0; i < len(scoreByCategoryList); i++ {
		if scoreByCategoryList[i] > 0 {
			finalScore += 1
		} else if scoreByCategoryList[i] < 0 {
			finalScore -= 1
		}
	}

	return finalScore, moneySupply

}

func main() {

	initMap()

	countriesName := []string{"united-states", "euro-area", "united-kingdom", "australia", "new-zealand", "canada", "switzerland", "japan"}

	for i := 0; i < len(countriesName); i++ {
		economicsData := fetchData(countriesName[i])
		economicScore, moneySupplyScore := calculateScore(economicsData)

		fmt.Printf("Country : %s\n", countriesName[i])
		fmt.Printf("Economic Data Score : %d\n", economicScore)
		fmt.Printf("Money Supply Score : %f\n", moneySupplyScore)
		fmt.Println()
	}

}
