package request

type ProductRequest struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Description string `json:"description"`
	CategoryId  int    `json:"category_id"`
}
