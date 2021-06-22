package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/wildangbudhi/forex-daily-analytics/services"
)

type Country struct {
	Economcis services.EconomcisData
	COT       services.CommitmentOfTrader
	Technical services.Technical
}

type PairData struct {
	EconomicsScore   int
	MoneySupplyScore float64
	COTScore         float64
	COTChangesScore  float64
	Technical        services.Technical
}

var (
	countriesName     []string                      = []string{"united-states", "euro-area", "united-kingdom", "australia", "new-zealand", "canada", "switzerland", "japan", "gold", "indonesia"}
	countriesCurrency []string                      = []string{"usd", "eur", "gbp", "aud", "nzd", "cad", "chf", "jpy", "xau", "idr"}
	technicalAllData  map[string]services.Technical = make(map[string]services.Technical)
	currencyMap       map[string]*Country           = make(map[string]*Country)
	pairs             [][]string                    = make([][]string, 0)
	pairData          map[string]*PairData          = make(map[string]*PairData)
)

func prepareData() {
	technicalAllData = services.NewTechnicalData()

	currencyMap = make(map[string]*Country)

	for i := 0; i < len(countriesName); i++ {
		var newCountry Country = Country{
			COT: services.NewCommitmentOfTrader(countriesName[i]),
		}

		if countriesCurrency[i] != "xau" {
			newCountry.Economcis = services.NewEconomicsData(countriesName[i])
		}

		currencyMap[countriesCurrency[i]] = &newCountry
	}

	var wg sync.WaitGroup

	for i := 0; i < len(countriesCurrency); i++ {

		if countriesCurrency[i] != "xau" {
			wg.Add(1)

			go func(currencyMap map[string]*Country, countryCurrency string, wg *sync.WaitGroup) {
				defer wg.Done()

				currencyMap[countryCurrency].Economcis.FetchData()
				currencyMap[countryCurrency].Economcis.CalculateScore()

			}(currencyMap, countriesCurrency[i], &wg)
		}

		wg.Add(1)

		go func(currencyMap map[string]*Country, countryCurrency string, wg *sync.WaitGroup) {
			defer wg.Done()

			currencyMap[countryCurrency].COT.FetchData()
			currencyMap[countryCurrency].COT.CalculateScore()

		}(currencyMap, countriesCurrency[i], &wg)

	}

	wg.Wait()

	pairs = [][]string{
		{"aud", "cad"},
		{"aud", "chf"},
		{"aud", "jpy"},
		{"aud", "nzd"},
		{"aud", "usd"},
		{"cad", "chf"},
		{"cad", "jpy"},
		{"chf", "jpy"},
		{"eur", "aud"},
		{"eur", "cad"},
		{"eur", "chf"},
		{"eur", "gbp"},
		{"eur", "jpy"},
		{"eur", "nzd"},
		{"eur", "usd"},
		{"gbp", "aud"},
		{"gbp", "cad"},
		{"gbp", "chf"},
		{"gbp", "jpy"},
		{"gbp", "nzd"},
		{"gbp", "usd"},
		{"nzd", "cad"},
		{"nzd", "chf"},
		{"nzd", "jpy"},
		{"nzd", "usd"},
		{"usd", "cad"},
		{"usd", "chf"},
		{"usd", "jpy"},
		{"xau", "usd"},
	}

	for i := 0; i < len(pairs); i++ {
		var baseCurrency, quoteCurrency string = strings.ToUpper(pairs[i][0]), strings.ToUpper(pairs[i][1])
		var pair PairData = PairData{}

		if pairs[i][0] != "xau" && pairs[i][1] != "xau" {
			pair.EconomicsScore = currencyMap[pairs[i][0]].Economcis.EconomicsScore - currencyMap[pairs[i][1]].Economcis.EconomicsScore
			pair.MoneySupplyScore = currencyMap[pairs[i][1]].Economcis.MoneySupply.Score - currencyMap[pairs[i][0]].Economcis.MoneySupply.Score
		} else {
			pair.EconomicsScore = currencyMap[pairs[i][1]].Economcis.EconomicsScore
			pair.MoneySupplyScore = currencyMap[pairs[i][1]].Economcis.MoneySupply.Score
		}

		pair.COTScore = currencyMap[pairs[i][0]].COT.Difference - currencyMap[pairs[i][1]].COT.Difference
		pair.COTChangesScore = currencyMap[pairs[i][0]].COT.ChangesDifferece - currencyMap[pairs[i][1]].COT.ChangesDifferece

		var pairDataKey string = fmt.Sprintf("%s%s", strings.ToLower(baseCurrency), strings.ToLower(quoteCurrency))

		pair.Technical = technicalAllData[pairDataKey]

		pairData[pairDataKey] = &pair

	}
}

func previewIndividuData() {

	fmt.Println()
	fmt.Println()

	fmt.Println("===================================================== INDIVIDUAL DATA =====================================================")
	fmt.Println()

	for i := 0; i < len(countriesCurrency); i++ {
		fmt.Printf("Country : %s\n", countriesName[i])
		fmt.Printf("Currency : %s\n", countriesCurrency[i])

		if countriesCurrency[i] != "xau" {
			fmt.Printf("Economic Data Score : %d\n", currencyMap[countriesCurrency[i]].Economcis.EconomicsScore)
			fmt.Printf(
				"Money Supply Score : (%.0f - %.0f) = %.6f%% \n",
				currencyMap[countriesCurrency[i]].Economcis.MoneySupply.Last,
				currencyMap[countriesCurrency[i]].Economcis.MoneySupply.Previous,
				currencyMap[countriesCurrency[i]].Economcis.MoneySupply.Score,
			)
		}

		fmt.Printf(
			"COT : (%.0f - %.0f) = %.0f \n",
			currencyMap[countriesCurrency[i]].COT.Long,
			currencyMap[countriesCurrency[i]].COT.Short,
			currencyMap[countriesCurrency[i]].COT.Difference,
		)

		fmt.Printf(
			"COT Changes: (%.2f - %.2f) = %.2f \n",
			currencyMap[countriesCurrency[i]].COT.LongChanges,
			currencyMap[countriesCurrency[i]].COT.ShortChanges,
			currencyMap[countriesCurrency[i]].COT.ChangesDifferece,
		)

		fmt.Println()
	}

	fmt.Println("===========================================================================================================================")

}

func previewPairData() {
	fmt.Println()
	fmt.Println()

	fmt.Println("===================================================== PAIRS DATA =====================================================")
	fmt.Println()

	for i := 0; i < len(pairs); i++ {

		var pairName string = fmt.Sprintf("%s%s", strings.ToLower(pairs[i][0]), strings.ToLower(pairs[i][1]))
		var pair *PairData = pairData[pairName]

		fmt.Printf("PAIR : %s\n", pairName)
		fmt.Printf("ECONOMICS SCORE : %d\n", pair.EconomicsScore)
		fmt.Printf("MONEY SUPPLY SCORE : %.6f\n", pair.MoneySupplyScore)
		fmt.Printf("COT SCORE : %.3f\n", pair.COTScore)
		fmt.Printf("COT CHANGES SCORE : %.3f\n", pair.COTChangesScore)
		fmt.Printf("TECHNICAL\n")
		fmt.Printf("Hourly : %s\n", pair.Technical.Hourly)
		fmt.Printf("Daily : %s\n", pair.Technical.Daily)
		fmt.Printf("Weekly : %s\n", pair.Technical.Weekly)
		fmt.Printf("Monthly : %s\n", pair.Technical.Monthly)

		fmt.Println()

	}

	fmt.Println("======================================================================================================================")
}

func main() {

	prepareData()
	previewIndividuData()
	previewPairData()

}
