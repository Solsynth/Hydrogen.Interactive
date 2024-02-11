package server

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"time"
)

func getOwnPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	id := c.Params("postId")
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	tx := database.C.Where(&models.Post{
		Alias:    id,
		AuthorID: user.ID,
	})

	post, err := services.GetPost(tx)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	tx = database.C.
		Where(&models.Post{ReplyID: &post.ID}).
		Where("published_at <= ? OR published_at IS NULL", time.Now()).
		Order("created_at desc")

	var count int64
	if err := tx.
		Model(&models.Post{}).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	posts, err := services.ListPost(tx, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"data":    post,
		"count":   count,
		"related": posts,
	})
}

func listOwnPost(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)

	user := c.Locals("principal").(models.Account)

	tx := database.C.
		Where(&models.Post{AuthorID: user.ID}).
		Where("published_at <= ? OR published_at IS NULL", time.Now()).
		Order("created_at desc")

	if realmId > 0 {
		tx = tx.Where(&models.Post{RealmID: lo.ToPtr(uint(realmId))})
	}

	if len(c.Query("category")) > 0 {
		tx = services.FilterPostWithCategory(tx, c.Query("category"))
	}

	if len(c.Query("tag")) > 0 {
		tx = services.FilterPostWithTag(tx, c.Query("tag"))
	}

	var count int64
	if err := tx.
		Model(&models.Post{}).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	posts, err := services.ListPost(tx, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  posts,
	})
}
