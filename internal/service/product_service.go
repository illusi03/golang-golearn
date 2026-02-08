package service

import (
	"context"

	"github.com/illusi03/golearn/internal/model"
	"github.com/illusi03/golearn/internal/repository"
)

type ProductService struct {
	productRepository *repository.ProductRepository
}

func NewProductService(productRepository *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (a *ProductService) FindAll(ctx context.Context, name string) ([]model.ProductModel, error) {
	return a.productRepository.FindAll(ctx, name)
}

func (a *ProductService) FindOne(ctx context.Context, id int) (*model.ProductModel, error) {
	return a.productRepository.FindOne(ctx, id)
}

func (a *ProductService) Delete(ctx context.Context, id int) (bool, error) {
	return a.productRepository.Delete(ctx, id)
}

func (a *ProductService) Update(ctx context.Context, model *model.ProductModel) (bool, error) {
	return a.productRepository.Update(ctx, model)
}

func (a *ProductService) Create(ctx context.Context, model *model.ProductModel) (*model.ProductModel, error) {
	return a.productRepository.Create(ctx, model)
}
