package service

import (
	"context"

	"github.com/illusi03/golearn/internal/model"
	"github.com/illusi03/golearn/internal/repository"
)

type CategoryService struct {
	categoryRepository *repository.CategoryRepository
}

func NewCategoryService(categoryRepository *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepository: categoryRepository,
	}
}

func (a *CategoryService) FindAll(ctx context.Context) ([]model.CategoryModel, error) {
	return a.categoryRepository.FindAll(ctx)
}

func (a *CategoryService) FindOne(ctx context.Context, id int) (*model.CategoryModel, error) {
	return a.categoryRepository.FindOne(ctx, id)
}

func (a *CategoryService) Delete(ctx context.Context, id int) (bool, error) {
	return a.categoryRepository.Delete(ctx, id)
}

func (a *CategoryService) Update(ctx context.Context, model *model.CategoryModel) (bool, error) {
	return a.categoryRepository.Update(ctx, model)
}

func (a *CategoryService) Create(ctx context.Context, model *model.CategoryModel) (*model.CategoryModel, error) {
	return a.categoryRepository.Create(ctx, model)
}
