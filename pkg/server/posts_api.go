package server

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"strings"
)

func listPost(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	var count int64
	var posts []*models.Post
	if err := database.C.
		Where(&models.Post{RealmID: nil}).
		Model(&models.Post{}).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := database.C.
		Where(&models.Post{RealmID: nil}).
		Limit(take).
		Offset(offset).
		Order("created_at desc").
		Preload("Author").
		Find(&posts).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	postIds := lo.Map(posts, func(item *models.Post, _ int) uint {
		return item.ID
	})

	var reactInfo []struct {
		PostID       uint
		LikeCount    int64
		DislikeCount int64
	}

	prefix := viper.GetString("database.prefix")
	database.C.Raw(fmt.Sprintf(`SELECT t.id                         as post_id,
       COALESCE(l.like_count, 0)    AS like_count,
       COALESCE(d.dislike_count, 0) AS dislike_count
FROM %sposts t
         LEFT JOIN (SELECT post_id, COUNT(*) AS like_count
                    FROM %spost_likes
                    GROUP BY post_id) l ON t.id = l.post_id
         LEFT JOIN (SELECT post_id, COUNT(*) AS dislike_count
                    FROM %spost_dislikes
                    GROUP BY post_id) d ON t.id = d.post_id
WHERE t.id IN (?)`, prefix, prefix, prefix), postIds).Scan(&reactInfo)

	postMap := lo.SliceToMap(posts, func(item *models.Post) (uint, *models.Post) {
		return item.ID, item
	})

	for _, info := range reactInfo {
		if post, ok := postMap[info.PostID]; ok {
			post.LikeCount = info.LikeCount
			post.DislikeCount = info.DislikeCount
		}
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  posts,
	})

}

func createPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Alias      string            `json:"alias"`
		Title      string            `json:"title"`
		Content    string            `json:"content" validate:"required"`
		Tags       []models.Tag      `json:"tags"`
		Categories []models.Category `json:"categories"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	} else if len(data.Alias) == 0 {
		data.Alias = strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	post, err := services.NewPost(user, data.Alias, data.Title, data.Content, data.Categories, data.Tags)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(post)
}

func reactPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("postId", 0)

	var post models.Post
	if err := database.C.Where(&models.Post{
		BaseModel: models.BaseModel{ID: uint(id)},
	}).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	switch strings.ToLower(c.Params("reactType")) {
	case "like":
		if positive, err := services.LikePost(user, post); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else {
			return c.SendStatus(lo.Ternary(positive, fiber.StatusCreated, fiber.StatusNoContent))
		}
	case "dislike":
		if positive, err := services.DislikePost(user, post); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else {
			return c.SendStatus(lo.Ternary(positive, fiber.StatusCreated, fiber.StatusNoContent))
		}
	default:
		return fiber.NewError(fiber.StatusBadRequest, "unsupported reaction")
	}
}
