package services

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CommitmentOfTrader struct {
	CountryCode      string
	CountryID        string
	Long             float64
	Short            float64
	Difference       float64
	LongChanges      float64
	ShortChanges     float64
	ChangesDifferece float64
}

var countryMap map[string]string = map[string]string{
	"united-states":  "098662",
	"euro-area":      "099741",
	"united-kingdom": "096742",
	"australia":      "232741",
	"new-zealand":    "112741",
	"canada":         "090741",
	"switzerland":    "092741",
	"japan":          "097741",
	"gold":           "088691",
}

func NewCommitmentOfTrader(countryCode string) CommitmentOfTrader {
	var object CommitmentOfTrader = CommitmentOfTrader{
		CountryCode: countryCode,
		CountryID:   countryMap[countryCode],
	}

	return object
}

func (obj *CommitmentOfTrader) FetchData() {

	fmt.Printf("[%s] Fetch COT Data\n", obj.CountryCode)

	res, err := http.Get("https://www.tradingster.com/cot/legacy-futures/" + obj.CountryID)

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

	doc.Find(".table").Children().Next().Children().Each(func(i int, row *goquery.Selection) {

		if i == 1 {
			row.Children().Each(func(j int, column *goquery.Selection) {

				dataString := strings.ReplaceAll(column.Text(), ",", "")

				if j == 0 {
					long, _ := strconv.ParseFloat(dataString, 64)
					obj.Long = long
					fmt.Printf("Long : %f\n", long)
				} else if j == 1 {
					short, _ := strconv.ParseFloat(dataString, 64)
					obj.Short = short
					fmt.Printf("Short : %f\n", short)
				}

			})
		} else if i == 3 {
			row.Children().Each(func(j int, column *goquery.Selection) {

				dataString := strings.ReplaceAll(column.Text(), ",", "")

				if j == 0 {
					longChanges, _ := strconv.ParseFloat(dataString, 64)
					obj.LongChanges = longChanges
				} else if j == 1 {
					shortChanges, _ := strconv.ParseFloat(dataString, 64)
					obj.ShortChanges = shortChanges
				}

			})
		}

	})

}

func (obj *CommitmentOfTrader) CalculateScore() {

	fmt.Printf("[%s] Calculate COT Score\n", obj.CountryCode)

	obj.Difference = obj.Long - obj.Short
	obj.ChangesDifferece = obj.LongChanges - obj.ShortChanges
}
