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
}

func main() {

	var countriesName []string = []string{"united-states", "euro-area", "united-kingdom", "australia", "new-zealand", "canada", "switzerland", "japan", "gold"}
	var countriesCurrency []string = []string{"usd", "eur", "gbp", "aud", "nzd", "cad", "chf", "jpy", "xau"}

	var cotData services.CommitmentOfTrader = services.NewCommitmentOfTrader(countriesName[5])
	cotData.FetchData()
	cotData.CalculateScore()

	var currencyMap map[string]*Country = make(map[string]*Country)

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

	var pairs [][]string = [][]string{
		{"aud", "cad"},
		{"aud", "chf"},
		{"aud", "jpy"},
		{"aud", "nzd"},
		{"aud", "usd"},
		{"cad", "chf"},
		{"cad", "jpy"},
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

	fmt.Println()
	fmt.Println()

	fmt.Println("===================================================== PAIRS DATA =====================================================")
	fmt.Println()

	for i := 0; i < len(pairs); i++ {

		var baseCurrency, quoteCurrency string = strings.ToUpper(pairs[i][0]), strings.ToUpper(pairs[i][1])
		var economicsScore int
		var moneySuppyScore, cotScore, cotChangesScore float64

		if pairs[i][0] != "xau" && pairs[i][1] != "xau" {
			economicsScore = currencyMap[pairs[i][0]].Economcis.EconomicsScore - currencyMap[pairs[i][1]].Economcis.EconomicsScore
			moneySuppyScore = currencyMap[pairs[i][0]].Economcis.MoneySupply.Score - currencyMap[pairs[i][1]].Economcis.MoneySupply.Score
		} else {
			economicsScore = currencyMap[pairs[i][1]].Economcis.EconomicsScore
			moneySuppyScore = currencyMap[pairs[i][1]].Economcis.MoneySupply.Score
		}

		cotScore = currencyMap[pairs[i][0]].COT.Difference - currencyMap[pairs[i][1]].COT.Difference
		cotChangesScore = currencyMap[pairs[i][0]].COT.ChangesDifferece - currencyMap[pairs[i][1]].COT.ChangesDifferece

		fmt.Printf("PAIR : %s%s\n", baseCurrency, quoteCurrency)
		fmt.Printf("ECONOMICS SCORE : %d\n", economicsScore)
		fmt.Printf("MONEY SUPPLY SCORE : %.6f\n", moneySuppyScore)
		fmt.Printf("COT SCORE : %.3f\n", cotScore)
		fmt.Printf("COT CHANGES SCORE : %.3f\n", cotChangesScore)

		fmt.Println()

	}

	fmt.Println("======================================================================================================================")

}
