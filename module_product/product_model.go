package module_product

type ProductModel struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

var LastProductId int = 0
var ProductDatas []*ProductModel
