package services

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Technical struct {
	Hourly  string
	Daily   string
	Weekly  string
	Monthly string
}

func fetchData(category string) []string {
	var url string

	if category == "currency" {
		url = "https://www.investing.com/currencies/service/Technical?pairid=0&sid=0.43038804868640823&smlID=1053843&category=Technical&download=true"
	} else if category == "commodity" {
		url = "https://www.investing.com/commodities/service/Technical?pairid=0&sid=0.43038804868640823&smlID=1053843&category=Technical&download=true"
	}

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)

	request.Header.Set("Cookie", "ses_id=MH4%2Bf2FuZW0%2Fe2BmYTAyMGMzMG5lYWZnNTxjZWFlZ3Fnc2ZoM2Q2cDM8OXcwMzUpNT0%2BN2FmN2w3MmVhMzBiNDA8PmxhNWVsP2hgbGE3MjVjYDBtZWFmYjVlY2JhYGdoZzVmNDNkNjIzbDlnMGw1bTUnPiJhJTcmN2VlNTNyYiUwPz5%2FYTJlbz89YGhhYzJkYzowbmUwZmY1YWNiYWVnf2cs;")
	request.Header.Set("User-Agent", "PostmanRuntime/7.28.0")

	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseString string = string(responseData)

	// fmt.Println(responseString)

	var responses []string = strings.Split(responseString, "\n")
	responses = responses[1 : len(responses)-1]

	return responses
}

func NewTechnicalData() map[string]Technical {

	var data []string = make([]string, 0)

	var currencyData []string = fetchData("currency")
	time.Sleep(1 * time.Second)
	var commodityData []string = fetchData("commodity")

	data = append(data, currencyData...)
	data = append(data, commodityData...)

	var technicalMap map[string]Technical = make(map[string]Technical)

	for i := 0; i < len(data); i++ {

		var temp string = data[i]
		var colTemp []string = strings.Split(temp, ",")

		for j := 0; j < len(colTemp); j++ {
			colTemp[j] = strings.ReplaceAll(colTemp[j], `"`, "")
		}

		colTemp[0] = strings.ReplaceAll(colTemp[0], "/", "")
		colTemp[0] = strings.ToLower(colTemp[0])

		if len(colTemp) > 5 {
			technicalMap[colTemp[0]] = Technical{
				Hourly:  colTemp[2],
				Daily:   colTemp[3],
				Weekly:  colTemp[4],
				Monthly: colTemp[5],
			}
		} else {
			technicalMap[colTemp[0]] = Technical{
				Hourly:  colTemp[1],
				Daily:   colTemp[2],
				Weekly:  colTemp[3],
				Monthly: colTemp[4],
			}
		}

	}

	return technicalMap

}
