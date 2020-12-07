package model

type Country struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Currency        string `json:"currency"`
	CurrencyCode    string `json:"currency_code"`
	CentralBank     string `json:"central_bank"`
	CentralBankCode string `json:"central_bank_code"`
}

type CountryRepository interface {
	Fetch() ([]Country, error)
}
