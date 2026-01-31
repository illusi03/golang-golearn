package model

type ProductModel struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Price        int    `json:"price"`
	CategoryID   *int   `json:"category_id"`
	CategoryName string `json:"category_name"`
}
