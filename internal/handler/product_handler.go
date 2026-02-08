package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/illusi03/golearn/internal/model"
	"github.com/illusi03/golearn/internal/request"
	"github.com/illusi03/golearn/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) Create(c fiber.Ctx) error {
	request := &request.ProductRequest{}
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Error validation json",
			"error":   err.Error(),
		})
	}

	productModel := &model.ProductModel{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}
	if request.CategoryId > 0 {
		productModel.CategoryID = &request.CategoryId
	}

	data, err := h.productService.Create(c.Context(), productModel)
	if err != nil {
		return fmt.Errorf("Error Occured : %w", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data created successfully",
		"data":    data,
	})
}

func (h *ProductHandler) Update(c fiber.Ctx) error {
	request := &request.ProductRequest{}
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Error validation json",
			"error":   err.Error(),
		})
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Error id product",
			"error":   nil,
		})
	}

	productModel := &model.ProductModel{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}
	if request.CategoryId > 0 {
		productModel.CategoryID = &request.CategoryId
	}

	data, err := h.productService.Update(c.Context(), productModel)
	if err != nil {
		return fmt.Errorf("Error Occured : %w", err)
	}

	if !data {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   data,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data updated successfully",
		"data":    data,
	})
}

func (h *ProductHandler) Delete(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Error id product",
			"error":   nil,
		})
	}

	data, err := h.productService.Delete(c.Context(), id)
	if err != nil {
		return fmt.Errorf("Error Occured : %w", err)
	}

	if !data {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   data,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data deleted successfully",
		"data":    data,
	})
}

func (h *ProductHandler) GetAll(c fiber.Ctx) error {
	name := c.Query("name")
	list, err := h.productService.FindAll(c.Context(), name)
	if err != nil {
		return fmt.Errorf("Error Occured : %w", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data fetched successfully",
		"data":    list,
	})
}

func (h *ProductHandler) GetDetail(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Error id product",
			"error":   nil,
		})
	}

	data, err := h.productService.FindOne(c.Context(), id)
	if err != nil {
		return fmt.Errorf("Error Occured: %w", err)
	}

	if data == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   data,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data fetched successfully",
		"data":    data,
	})
}
