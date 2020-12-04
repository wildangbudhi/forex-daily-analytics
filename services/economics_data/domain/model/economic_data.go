package model

type Country struct {
	ID              string `gorm:"primaryKey"`
	Name            string
	CurrencyCode    string
	Curremcy        string
	CentralBank     string
	CentralBankCode string
}

type EconomicData struct {
	CountryID string
	Category  string
	Indicator string
	Last      float64
	Previous  float64
	Country   Country `gorm:"foreignKey:CountryID"`
}
