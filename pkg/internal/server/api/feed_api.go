package api

import (
	"fmt"
	"time"

	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

type FeedRecord struct {
	Type      string         `json:"type"`
	Data      map[string]any `json:"data"`
	CreatedAt time.Time      `json:"created_at"`
}

func listFeed(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)

	postTx := services.FilterPostDraft(database.C)

	if realmId > 0 {
		if realm, err := services.GetRealmWithExtID(uint(realmId)); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("realm was not found: %v", err))
		} else {
			postTx = services.FilterPostWithRealm(postTx, realm.ID)
		}
	}

	if len(c.Query("authorId")) > 0 {
		var author models.Account
		if err := database.C.Where(&models.Account{Name: c.Query("authorId")}).First(&author).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		postTx = postTx.Where("author_id = ?", author.ID)
	}

	if len(c.Query("category")) > 0 {
		postTx = services.FilterPostWithCategory(postTx, c.Query("category"))
	}
	if len(c.Query("tag")) > 0 {
		postTx = services.FilterPostWithTag(postTx, c.Query("tag"))
	}

	postCountTx := postTx

	postCount, err := services.CountPost(postCountTx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	postItems, err := services.ListPost(postTx, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var feed []FeedRecord

	encodeToFeed := func(t string, in any, createdAt time.Time) FeedRecord {
		var result map[string]any
		raw, _ := jsoniter.Marshal(in)
		_ = jsoniter.Unmarshal(raw, &result)

		return FeedRecord{
			Type:      t,
			Data:      result,
			CreatedAt: createdAt,
		}
	}

	for _, post := range postItems {
		feed = append(feed, encodeToFeed("post", post, post.CreatedAt))
	}

	return c.JSON(fiber.Map{
		"count": postCount,
		"data":  feed,
	})
}

func listDraftMixed(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	postTx := services.FilterPostWithAuthorDraft(database.C, user.ID)
	postCountTx := postTx

	postCount, err := services.CountPost(postCountTx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	postItems, err := services.ListPost(postTx, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var feed []FeedRecord

	encodeToFeed := func(t string, in any, createdAt time.Time) FeedRecord {
		var result map[string]any
		raw, _ := jsoniter.Marshal(in)
		_ = jsoniter.Unmarshal(raw, &result)

		return FeedRecord{
			Type:      t,
			Data:      result,
			CreatedAt: createdAt,
		}
	}

	for _, post := range postItems {
		feed = append(feed, encodeToFeed("post", post, post.CreatedAt))
	}

	return c.JSON(fiber.Map{
		"count": postCount,
		"data":  feed,
	})
}
