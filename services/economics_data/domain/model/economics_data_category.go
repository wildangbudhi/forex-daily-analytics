package model

type EconomicsDataCategory struct {
	Id           int    `json:"id"`
	CategoryName string `json:"category_name"`
}

type EconomicsDataCategoryRepository interface {
	GetByCategoryName(categoryName string) (*EconomicsDataCategory, error)
	Insert(categoryName string) (int, error)
}
