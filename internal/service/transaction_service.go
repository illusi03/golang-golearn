package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/illusi03/golearn/internal/model"
	"github.com/illusi03/golearn/internal/repository"
	"github.com/illusi03/golearn/internal/request"
)

type TransactionService struct {
	transactionRepository *repository.TransactionRepository
	productRepository     *repository.ProductRepository
}

func NewTransactionService(
	transactionRepository *repository.TransactionRepository,
	productRepository *repository.ProductRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		productRepository:     productRepository,
	}
}

func (s *TransactionService) Checkout(
	ctx context.Context,
	req *request.CheckoutRequest,
) (*model.TransactionModel, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("checkout items cannot be empty")
	}

	seenProducts := make(map[int]bool)
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than 0")
		}
		if seenProducts[item.ProductID] {
			return nil, errors.New("the requested product was duplicate")
		}
		seenProducts[item.ProductID] = true
	}

	var totalAmount int
	var details []model.TransactionDetailModel

	for _, item := range req.Items {
		productID := item.ProductID
		quantity := item.Quantity
		product, err := s.productRepository.FindOne(ctx, productID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, fmt.Errorf("product with id %d not found", productID)
		}

		if product.Stock < quantity {
			return nil, fmt.Errorf("insufficient stock for product %s (available: %d, requested: %d)",
				product.Name, product.Stock, quantity)
		}

		subtotal := product.Price * quantity
		totalAmount += subtotal

		details = append(details, model.TransactionDetailModel{
			ProductID:   productID,
			ProductName: &product.Name,
			Quantity:    quantity,
			Subtotal:    subtotal,
		})
	}

	transaction := &model.TransactionModel{
		TotalAmount: totalAmount,
		Details:     details,
	}

	return s.transactionRepository.Create(ctx, transaction)
}

func (s *TransactionService) FindAll(
	ctx context.Context,
) ([]model.TransactionModel, error) {
	return s.transactionRepository.FindAll(ctx)
}
