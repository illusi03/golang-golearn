package module_product

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type ProductHandler struct {
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (h *ProductHandler) Create(c fiber.Ctx) error {
	request := &ProductRequest{}
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Error validation json",
			"error":   err.Error(),
		})
	}

	LastProductId = LastProductId + 1
	newData := &ProductModel{
		ID:          LastProductId,
		Name:        request.Name,
		Price:       request.Price,
		Description: request.Description,
	}
	ProductDatas = append(ProductDatas, newData)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data created successfully",
		"data":    newData,
	})
}

func (h *ProductHandler) Update(c fiber.Ctx) error {
	request := &ProductRequest{}
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

	var productModel *ProductModel
	for i, v := range ProductDatas {
		if v.ID == id {
			ProductDatas[i].Name = request.Name
			ProductDatas[i].Price = request.Price
			ProductDatas[i].Description = request.Description
			productModel = v
		}
	}

	if productModel == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   productModel,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data updated successfully",
		"data":    productModel,
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

	var productModel *ProductModel
	for i, v := range ProductDatas {
		if ProductDatas[i].ID == id {
			ProductDatas = append(ProductDatas[:i], ProductDatas[i+1:]...)
			productModel = v
			break
		}
	}

	if productModel == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   productModel,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data deleted successfully",
		"data":    productModel,
	})
}

func (h *ProductHandler) GetAll(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data fetched successfully",
		"data":    ProductDatas,
	})
}
