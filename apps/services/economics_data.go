package services

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

type MoneySupplyData struct {
	Last     float64
	Previous float64
	Score    float64
}

type EconomcisData struct {
	CountryCode    string
	EconomicsScore int
	MoneySupply    MoneySupplyData
	economicsData  []tradingEconomicsData
	mapScore       map[string]map[string]func(last, prev float64) int
}

func NewEconomicsData(countryCode string) EconomcisData {
	var object EconomcisData = EconomcisData{CountryCode: countryCode}
	object.initMapScore()

	return object
}

func (obj *EconomcisData) FetchData() {

	fmt.Printf("[%s] Fetch Economics Data\n", obj.CountryCode)

	res, err := http.Get("https://tradingeconomics.com/" + obj.CountryCode + "/indicators")

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

	fmt.Println(tradingEconomicData)

	obj.economicsData = tradingEconomicData
}

func (obj *EconomcisData) greaterBetter(last, prev float64) int {
	if last > prev {
		return 1
	} else if last == prev {
		return 0
	} else {
		return -1
	}
}

func (obj *EconomcisData) lessBetter(last, prev float64) int {
	if last < prev {
		return 1
	} else if last == prev {
		return 0
	} else {
		return -1
	}
}

func (obj *EconomcisData) initMapScore() {

	obj.mapScore = make(map[string]map[string]func(last float64, prev float64) int)

	obj.mapScore["GDP"] = map[string]func(last float64, prev float64) int{
		"GDP Growth Rate":        obj.greaterBetter,
		"GDP Annual Growth Rate": obj.greaterBetter,
		"GDP Growth Annualized":  obj.greaterBetter,
	}

	obj.mapScore["Labour"] = map[string]func(last float64, prev float64) int{
		"Unemployment Rate": obj.lessBetter,
		"Non Farm Payrolls": obj.greaterBetter,
	}

	obj.mapScore["Prices"] = map[string]func(last float64, prev float64) int{
		"Inflation Rate":     obj.greaterBetter,
		"Inflation Rate Mom": obj.greaterBetter,
	}

	obj.mapScore["Trade"] = map[string]func(last float64, prev float64) int{
		"Balance of Trade":       obj.greaterBetter,
		"Current Account":        obj.greaterBetter,
		"Current Account to GDP": obj.greaterBetter,
	}

	obj.mapScore["Government"] = map[string]func(last float64, prev float64) int{
		"Government Debt to GDP": obj.lessBetter,
		"Government Budget":      obj.greaterBetter,
	}

	obj.mapScore["Business"] = map[string]func(last float64, prev float64) int{
		"Business Confidence":   obj.greaterBetter,
		"Manufacturing PMI":     obj.greaterBetter,
		"Non Manufacturing PMI": obj.greaterBetter,
		"Services PMI":          obj.greaterBetter,
	}

	obj.mapScore["Consumer"] = map[string]func(last float64, prev float64) int{
		"Consumer Confidence": obj.greaterBetter,
		"Retail Sales MoM":    obj.greaterBetter,
	}

	obj.mapScore["Housing"] = map[string]func(last float64, prev float64) int{
		"Building Permits": obj.greaterBetter,
	}

	obj.mapScore["Taxes"] = map[string]func(last float64, prev float64) int{
		"Corporate Tax Rate":       obj.greaterBetter,
		"Personal Income Tax Rate": obj.greaterBetter,
	}

}

func (obj *EconomcisData) CalculateScore() {

	fmt.Printf("[%s] Calculate Economics Score and Money Supply Growth\n", obj.CountryCode)

	var score map[string]map[string]int = make(map[string]map[string]int)
	var moneySupply MoneySupplyData = MoneySupplyData{}

	for i := 0; i < len(obj.economicsData); i++ {
		var category, indicator string = obj.economicsData[i].Category, obj.economicsData[i].Indicator
		var last, prev float64 = obj.economicsData[i].Last, obj.economicsData[i].Previous

		if category == "Money" && indicator == "Money Supply M1" {
			moneySupply.Last = last
			moneySupply.Previous = prev
			moneySupply.Score = (last - prev) / prev
		} else {

			indicatorValue, isIndicatorValueExist := obj.mapScore[category]

			if isIndicatorValueExist {
				isBetterFunction, ok := indicatorValue[indicator]

				if ok {

					if len(score[category]) == 0 {
						score[category] = make(map[string]int)
					}

					score[category][indicator] = isBetterFunction(last, prev)
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

	obj.EconomicsScore = finalScore
	obj.MoneySupply = moneySupply

}
