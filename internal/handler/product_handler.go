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

	data, err := h.productService.Create(c, &model.ProductModel{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
	})
	if err != nil {
		return fmt.Errorf("Terdapat error : %w", err)
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

	data, err := h.productService.Update(c, &model.ProductModel{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
	})
	if err != nil {
		return fmt.Errorf("Terdapat error : %w", err)
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

	data, err := h.productService.Delete(c, id)
	if err != nil {
		return fmt.Errorf("Terdapat error : %w", err)
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
	list, err := h.productService.FindAll(c)
	if err != nil {
		return fmt.Errorf("Terdapat error : %w", err)
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

	data, err := h.productService.FindOne(c, id)
	if err != nil {
		return fmt.Errorf("Terdapat error: %w", err)
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
