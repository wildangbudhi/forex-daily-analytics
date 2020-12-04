package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type EconomicData struct {
	Category  string
	Indicator string
	Last      float64
	Previous  float64
}

type Rules struct {
	CategoryIndex int
	Function      func(float64, float64) int
}

func fetchData(countryCode string) []EconomicData {

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

	economicData := make([]EconomicData, 0)
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

							row := EconomicData{Category: category}

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

							economicData = append(economicData, row)

						})
					}
				})

			}
		})

	}

	return economicData
}

func higherLastBetterFunction(previous, last float64) int {
	if last > previous {
		return 1
	} else if last == previous {
		return 0
	} else {
		return -1
	}
}

func lowerLastBetterFunction(previous, last float64) int {
	if last < previous {
		return 1
	} else if last == previous {
		return 0
	} else {
		return -1
	}
}

func inflationFunction(previous, last float64) int {

	if last == 2.0 {
		return 1
	}

	previousDistance := math.Abs(previous - 2.0)
	lastDistance := math.Abs(last - 2.0)

	if (last > 2.0 && previous < 2.0) || (previous > 2.0 && last < 2.0) {
		return -1
	}

	if lastDistance < previousDistance {
		return 1
	} else if lastDistance == previousDistance {
		return 0
	} else {
		return -1
	}

}

func calculateScore(economicData []EconomicData) (float64, []float64) {

	categoryMap := map[string]Rules{
		"GDP Growth Rate": Rules{
			CategoryIndex: 0,
			Function:      higherLastBetterFunction,
		},
		"GDP Annual Growth Rate": Rules{
			CategoryIndex: 0,
			Function:      higherLastBetterFunction,
		},
		"Unemployment Rate": Rules{
			CategoryIndex: 1,
			Function:      lowerLastBetterFunction,
		},
		"Non Farm Payrolls": Rules{
			CategoryIndex: 1,
			Function:      higherLastBetterFunction,
		},
		"Inflation Rate": Rules{
			CategoryIndex: 2,
			Function:      inflationFunction,
		},
		"Inflation Rate Mom": Rules{
			CategoryIndex: 2,
			Function:      inflationFunction,
		},
		"Balance of Trade": Rules{
			CategoryIndex: 3,
			Function:      higherLastBetterFunction,
		},
		"Current Account": Rules{
			CategoryIndex: 3,
			Function:      higherLastBetterFunction,
		},
		"Current Account to GDP": Rules{
			CategoryIndex: 3,
			Function:      higherLastBetterFunction,
		},
		"Government Debt to GDP": Rules{
			CategoryIndex: 4,
			Function:      lowerLastBetterFunction,
		},
		"Government Budget": Rules{
			CategoryIndex: 4,
			Function:      higherLastBetterFunction,
		},
		"Business Confidence": Rules{
			CategoryIndex: 5,
			Function:      higherLastBetterFunction,
		},
		"Manufacturing PMI": Rules{
			CategoryIndex: 5,
			Function:      higherLastBetterFunction,
		},
		"Non Manufacturing PMI": Rules{
			CategoryIndex: 5,
			Function:      higherLastBetterFunction,
		},
		"Services PMI": Rules{
			CategoryIndex: 5,
			Function:      higherLastBetterFunction,
		},
		"Consumer Confidence": Rules{
			CategoryIndex: 6,
			Function:      higherLastBetterFunction,
		},
		"Retail Sales MoM": Rules{
			CategoryIndex: 6,
			Function:      higherLastBetterFunction,
		},
		"Building Permits": Rules{
			CategoryIndex: 7,
			Function:      higherLastBetterFunction,
		},
		"Corporate Tax Rate": Rules{
			CategoryIndex: 8,
			Function:      higherLastBetterFunction,
		},
		"Personal Income Tax Rate": Rules{
			CategoryIndex: 8,
			Function:      higherLastBetterFunction,
		},
	}

	scoreList := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0}

	for _, ed := range economicData {
		val, ok := categoryMap[ed.Indicator]

		if ok {
			conditionValue := val.Function(ed.Previous, ed.Last)
			scoreList[categoryMap[ed.Indicator].CategoryIndex] += float64(conditionValue)
		}
	}

	score := float64(0.0)

	for i := 0; i < len(scoreList); i++ {
		if scoreList[i] > 0.0 {
			score += 1.0
		} else if scoreList[i] < 0.0 {
			score -= 1.0
		}
	}

	return score, scoreList
}

func main() {

	// var wg sync.WaitGroup
	countryCodes := []string{"euro-area", "united-states", "united-kingdom", "canada", "switzerland", "japan", "australia", "new-zealand"}

	// for _, countryCode := range countryCodes {
	// 	wg.Add(1)
	// 	go func(cc string) {
	// 		defer wg.Done()
	// 		data := fetchData(cc)
	// 		score := calculateScore(data)
	// 		fmt.Printf("%s --> %f\n", cc, score)
	// 	}(countryCode)
	// }

	// wg.Wait()

	for _, countryCode := range countryCodes {
		data := fetchData(countryCode)
		score, scoreList := calculateScore(data)
		fmt.Printf("%s --> %f --> ", countryCode, score)
		fmt.Println(scoreList)
	}

	// data := fetchData("united-states")

	// for i := 0; i < len(data); i++ {
	// 	fmt.Println(data[i])
	// }
}
