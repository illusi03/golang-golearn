package module_category

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type CategoryHandler struct {
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

func (h *CategoryHandler) Create(c fiber.Ctx) error {
	request := &CategoryRequest{}
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Error validation json",
			"error":   err.Error(),
		})
	}

	LastCategoryId = LastCategoryId + 1
	newData := &CategoryModel{
		ID:          LastCategoryId,
		Name:        request.Name,
		Description: request.Description,
	}
	CategoryDatas = append(CategoryDatas, newData)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data created successfully",
		"data":    newData,
	})
}

func (h *CategoryHandler) Update(c fiber.Ctx) error {
	request := &CategoryRequest{}
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

	var productModel *CategoryModel
	for i, v := range CategoryDatas {
		if v.ID == id {
			CategoryDatas[i].Name = request.Name
			CategoryDatas[i].Description = request.Description
			productModel = v
		}
	}

	if productModel == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
			"error":   productModel,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data updated successfully",
		"data":    productModel,
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

	var productModel *CategoryModel
	for i, v := range CategoryDatas {
		if CategoryDatas[i].ID == id {
			CategoryDatas = append(CategoryDatas[:i], CategoryDatas[i+1:]...)
			productModel = v
			break
		}
	}

	if productModel == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
			"error":   productModel,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data deleted successfully",
		"data":    productModel,
	})
}

func (h *CategoryHandler) GetAll(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data fetched successfully",
		"data":    CategoryDatas,
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

	var productModel *CategoryModel
	for _, v := range CategoryDatas {
		if v.ID == id {
			productModel = v
		}
	}

	if productModel == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
			"error":   productModel,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data fetched successfully",
		"data":    productModel,
	})
}
