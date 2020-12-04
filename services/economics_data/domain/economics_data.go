package domain

type Rules struct {
	Category      string
	CategoryIndex int
	Function      func(float64, float64) int
}
