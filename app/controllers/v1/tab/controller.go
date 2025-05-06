package tab_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"hintword.com/api/app/common/validator"
	"hintword.com/api/app/database"
	"hintword.com/api/app/models"
	userService "hintword.com/api/app/services/user"
	"net/http"
)

func CreateUpdateCollection(c *fiber.Ctx) error {
	params := PayloadCollection{}
	if err := validator.ParseBodyAndValidate(c, &params); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  err,
		})
	}

	userDetails := userService.GetUserObject(c)
	collection := models.Collections{}
	collection, _ = collection.FindById(params.CollectionID)
	if collection.CollectionID == "" && params.CollectionID != "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  "invalid collection id given",
		})
	}

	if params.CollectionID == "" {
		collection.UserID = userDetails.UserId
	}
	collection.Name = params.Name
	collection.Status = params.Status
	collection, err := collection.Save()
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -3,
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": 1,
		"data":   collection,
	})
}

func CreateUpdateTab(c *fiber.Ctx) error {
	params := PayloadTab{}
	if err := validator.ParseBodyAndValidate(c, &params); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  err,
		})
	}
	userDetails := userService.GetUserObject(c)
	tab := models.Tabs{}
	tab, _ = tab.FindById(params.TabID)
	if tab.TabID == "" && params.TabID != "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  "invalid tab id given",
		})
	}

	collection := models.Collections{}
	collection, _ = collection.FindById(params.CollectionID)
	if collection.CollectionID == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -2,
			"error":  "invalid collection id given",
		})
	}

	if tab.TabID == "" {
		tab.UserID = userDetails.UserId
	}
	tab.Title = params.Title
	tab.URL = params.URL
	tab.FaviconURL = params.FaviconURL
	tab.Status = params.Status
	tab.Sequence = params.Sequence
	tab.CollectionID = collection.CollectionID
	tab, err := tab.Save()
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -3,
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": 1,
		"data":   tab,
	})
}

func GetCollection(c *fiber.Ctx) error {
	userDetails := userService.GetUserObject(c)
	var collections []models.Collections

	if err := database.MysqlDB.
		Where("user_id = ?", userDetails.UserId).
		Where("status = ?", models.StatusActive).
		Order("created_at desc").
		Find(&collections).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": -1,
			"error":  fmt.Sprintf("Failed to fetch notes: %v", err),
		})
	}

	for i, collection := range collections {
		tab := models.Tabs{}
		tabs, _ := tab.FindByCollectionId(collection.CollectionID)
		collections[i].Tabs = tabs
	}

	return c.JSON(fiber.Map{
		"status":     1,
		"collection": collections,
	})
}
