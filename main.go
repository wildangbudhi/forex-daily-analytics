package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type EconomicData struct {
	Key      string
	Last     float64
	Previous float64
}

type Rules struct {
	Category      string
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

	var temp *goquery.Selection

	doc.Find("#overview").Children().Each(func(i int, s *goquery.Selection) {
		if i == 3 {
			temp = s
		}
	})

	temp.Children().Children().Each(func(i int, s *goquery.Selection) {
		if i == 1 {
			temp = s
		}
	})

	economicData := make([]EconomicData, 0)

	temp.Children().Each(func(i int, s *goquery.Selection) {

		var row EconomicData

		s.Children().Each(func(j int, komp *goquery.Selection) {

			strKomp := strings.Replace(komp.Text(), "\n", "", -1)
			strKomp = regexp.MustCompile(`(\(.*\))`).ReplaceAllLiteralString(strKomp, "")
			strKomp = strings.TrimSpace(strKomp)

			if j == 0 {
				row.Key = strKomp
			} else if j == 1 {
				komponen, err := strconv.ParseFloat(strKomp, 64)

				if err != nil {
					log.Fatal(err)
				}

				row.Last = float64(komponen)
			} else if j == 3 {
				komponen, err := strconv.ParseFloat(strKomp, 64)

				if err != nil {
					log.Fatal(err)
				}

				row.Previous = float64(komponen)
			}
		})

		economicData = append(economicData, row)

	})

	return economicData
}

func calculateScore(economicData []EconomicData) float64 {

	categoryMap := map[string]Rules{
		"GDP Growth Rate": Rules{
			Category:      "GDP",
			CategoryIndex: 0,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"GDP Annual Growth Rate": Rules{
			Category:      "GDP",
			CategoryIndex: 0,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Unemployment Rate": Rules{
			Category:      "Labour",
			CategoryIndex: 1,
			Function: func(previous, last float64) int {
				if last < previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Non Farm Payrolls": Rules{
			Category:      "Labour",
			CategoryIndex: 1,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Inflation Rate": Rules{
			Category:      "Prices",
			CategoryIndex: 2,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Inflation Rate Mom": Rules{
			Category:      "Prices",
			CategoryIndex: 2,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Balance of Trade": Rules{
			Category:      "Trade",
			CategoryIndex: 3,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Current Account": Rules{
			Category:      "Trade",
			CategoryIndex: 3,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Current Account to GDP": Rules{
			Category:      "Trade",
			CategoryIndex: 3,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Government Debt to GDP": Rules{
			Category:      "Government",
			CategoryIndex: 4,
			Function: func(previous, last float64) int {
				if last < previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Government Budget": Rules{
			Category:      "Government",
			CategoryIndex: 4,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Business Confidence": Rules{
			Category:      "Business",
			CategoryIndex: 5,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Manufacturing PMI": Rules{
			Category:      "Business",
			CategoryIndex: 5,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Non Manufacturing PMI": Rules{
			Category:      "Business",
			CategoryIndex: 5,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Services PMI": Rules{
			Category:      "Business",
			CategoryIndex: 5,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Consumer Confidence": Rules{
			Category:      "Consumer",
			CategoryIndex: 6,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Retail Sales MoM": Rules{
			Category:      "Consumer",
			CategoryIndex: 6,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Building Permits": Rules{
			Category:      "Housing",
			CategoryIndex: 7,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Corporate Tax Rate": Rules{
			Category:      "Taxes",
			CategoryIndex: 8,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
		"Personal Income Tax Rate": Rules{
			Category:      "Taxes",
			CategoryIndex: 8,
			Function: func(previous, last float64) int {
				if last > previous {
					return 1
				} else if last == previous {
					return 0
				} else {
					return -1
				}
			},
		},
	}

	scoreList := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0}

	for _, ed := range economicData {
		val, ok := categoryMap[ed.Key]

		if ok {
			conditionValue := val.Function(ed.Previous, ed.Last)
			scoreList[categoryMap[ed.Key].CategoryIndex] += float64(conditionValue)
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

	return score
}

func main() {

	var wg sync.WaitGroup
	countryCodes := []string{"euro-area", "united-states", "united-kingdom", "canada", "switzerland", "japan", "australia", "new-zealand"}

	for _, countryCode := range countryCodes {
		wg.Add(1)
		go func(cc string) {
			defer wg.Done()
			data := fetchData(cc)
			score := calculateScore(data)
			fmt.Printf("%s --> %f\n", cc, score)
		}(countryCode)
	}

	wg.Wait()

}
