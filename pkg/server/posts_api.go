package server

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
)

func listPost(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	var count int64
	var posts []models.Post
	if err := database.C.
		Model(&models.Post{}).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := database.C.
		Limit(take).
		Offset(offset).
		Order("created_at desc").
		Preload("Author").
		Find(&posts).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  posts,
	})

}

func createPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Alias      string            `json:"alias" validate:"required"`
		Title      string            `json:"title" validate:"required"`
		Content    string            `json:"content" validate:"required"`
		Tags       []models.Tag      `json:"tags"`
		Categories []models.Category `json:"categories"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	post, err := services.NewPost(user, data.Alias, data.Title, data.Content, data.Categories, data.Tags)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(post)
}
