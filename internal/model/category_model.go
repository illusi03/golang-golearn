package model

type CategoryModel struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var LastCategoryId int = 0
var CategoryDatas []*CategoryModel
