package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/illusi03/golearn/internal/model"
	"github.com/illusi03/golearn/internal/request"
	"github.com/illusi03/golearn/internal/service"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) Create(c fiber.Ctx) error {
	request := &request.CategoryRequest{}
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Error validation json",
			"error":   err.Error(),
		})
	}

	data, err := h.categoryService.Create(c, &model.CategoryModel{
		Name:        request.Name,
		Description: request.Description,
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

func (h *CategoryHandler) Update(c fiber.Ctx) error {
	request := &request.CategoryRequest{}
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
			"message": "Error id category",
			"error":   nil,
		})
	}

	data, err := h.categoryService.Update(c, &model.CategoryModel{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
	})
	if err != nil {
		return fmt.Errorf("Terdapat error : %w", err)
	}

	if !data {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
			"error":   data,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data updated successfully",
		"data":    data,
	})
}

func (h *CategoryHandler) Delete(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Error id category",
			"error":   nil,
		})
	}

	data, err := h.categoryService.Delete(c, id)
	if err != nil {
		return fmt.Errorf("Terdapat error : %w", err)
	}

	if !data {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
			"error":   data,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data deleted successfully",
		"data":    data,
	})
}

func (h *CategoryHandler) GetAll(c fiber.Ctx) error {
	list, err := h.categoryService.FindAll(c)
	if err != nil {
		return fmt.Errorf("Terdapat error : %w", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data fetched successfully",
		"data":    list,
	})
}

func (h *CategoryHandler) GetDetail(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Error id category",
			"error":   nil,
		})
	}

	data, err := h.categoryService.FindOne(c, id)
	if err != nil {
		return fmt.Errorf("Terdapat error: %w", err)
	}

	if data == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
			"error":   data,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data fetched successfully",
		"data":    data,
	})
}
