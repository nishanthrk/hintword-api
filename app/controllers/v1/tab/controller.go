package tab_controller

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"hintword.com/api/app/common/validator"
	"hintword.com/api/app/models"
	userService "hintword.com/api/app/services/user"
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

	// Get collections ordered by sequence
	collectionsModel := models.Collections{}
	collections, err := collectionsModel.FindByUserOrderedBySequence(userDetails.UserId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": -1,
			"error":  fmt.Sprintf("Failed to fetch collections: %v", err),
		})
	}

	// Load tabs for each collection
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

func ReorderCollections(c *fiber.Ctx) error {
	params := PayloadReorderCollections{}
	if err := validator.ParseBodyAndValidate(c, &params); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  err,
		})
	}

	userDetails := userService.GetUserObject(c)

	// Verify all collections belong to the user
	collectionsModel := models.Collections{}
	for _, collectionUpdate := range params.Collections {
		collection, err := collectionsModel.FindById(collectionUpdate.CollectionID)
		if err != nil || collection.CollectionID == "" {
			return c.Status(http.StatusNotFound).JSON(&fiber.Map{
				"status": -1,
				"error":  fmt.Sprintf("Collection %s not found", collectionUpdate.CollectionID),
			})
		}

		if collection.UserID != userDetails.UserId {
			return c.Status(http.StatusForbidden).JSON(&fiber.Map{
				"status": -1,
				"error":  fmt.Sprintf("Access denied to collection %s", collectionUpdate.CollectionID),
			})
		}
	}

	// Convert to the format expected by BulkUpdateSequences
	updates := make([]struct {
		CollectionID string `json:"collection_id"`
		Sequence     int64  `json:"sequence"`
	}, len(params.Collections))

	for i, collectionUpdate := range params.Collections {
		updates[i] = struct {
			CollectionID string `json:"collection_id"`
			Sequence     int64  `json:"sequence"`
		}{
			CollectionID: collectionUpdate.CollectionID,
			Sequence:     collectionUpdate.Sequence,
		}
	}

	// Update sequences in bulk
	err := collectionsModel.BulkUpdateSequences(updates)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"status": -1,
			"error":  fmt.Sprintf("Failed to reorder collections: %v", err),
		})
	}

	// Return updated collections in new order
	collections, err := collectionsModel.FindByUserOrderedBySequence(userDetails.UserId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"status": -1,
			"error":  fmt.Sprintf("Failed to fetch reordered collections: %v", err),
		})
	}

	// Load tabs for each collection
	for i, collection := range collections {
		tab := models.Tabs{}
		tabs, _ := tab.FindByCollectionId(collection.CollectionID)
		collections[i].Tabs = tabs
	}

	return c.JSON(fiber.Map{
		"status":     1,
		"message":    "Collections reordered successfully",
		"collection": collections,
	})
}

func ReorderTab(c *fiber.Ctx) error {
	params := PayloadReorderTabs{}
	if err := validator.ParseBodyAndValidate(c, &params); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  err,
		})
	}

	userDetails := userService.GetUserObject(c)

	// Verify collection belongs to the user
	collectionsModel := models.Collections{}
	collection, err := collectionsModel.FindById(params.CollectionID)
	if err != nil || collection.CollectionID == "" {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"status": -1,
			"error":  "Collection not found",
		})
	}

	if collection.UserID != userDetails.UserId {
		return c.Status(http.StatusForbidden).JSON(&fiber.Map{
			"status": -1,
			"error":  "Access denied to this collection",
		})
	}

	// Verify all tabs belong to the user and collection
	tabModel := models.Tabs{}
	for _, tabUpdate := range params.Tabs {
		tab, err := tabModel.FindById(tabUpdate.TabID)
		if err != nil || tab.TabID == "" {
			return c.Status(http.StatusNotFound).JSON(&fiber.Map{
				"status": -1,
				"error":  fmt.Sprintf("Tab %s not found", tabUpdate.TabID),
			})
		}

		if tab.UserID != userDetails.UserId {
			return c.Status(http.StatusForbidden).JSON(&fiber.Map{
				"status": -1,
				"error":  fmt.Sprintf("Access denied to tab %s", tabUpdate.TabID),
			})
		}

		if tab.CollectionID != params.CollectionID {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"status": -1,
				"error":  fmt.Sprintf("Tab %s does not belong to collection %s", tabUpdate.TabID, params.CollectionID),
			})
		}
	}

	// Convert to the format expected by BulkUpdateSequences
	updates := make([]struct {
		TabID    string `json:"tab_id"`
		Sequence int64  `json:"sequence"`
	}, len(params.Tabs))

	for i, tabUpdate := range params.Tabs {
		updates[i] = struct {
			TabID    string `json:"tab_id"`
			Sequence int64  `json:"sequence"`
		}{
			TabID:    tabUpdate.TabID,
			Sequence: tabUpdate.Sequence,
		}
	}

	// Update sequences in bulk
	err = tabModel.BulkUpdateSequences(updates)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"status": -1,
			"error":  fmt.Sprintf("Failed to reorder tabs: %v", err),
		})
	}

	// Return updated tabs in new order
	tabs, err := tabModel.FindByCollectionId(params.CollectionID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"status": -1,
			"error":  fmt.Sprintf("Failed to fetch reordered tabs: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"status":  1,
		"message": "Tabs reordered successfully",
		"tabs":    tabs,
	})
}
