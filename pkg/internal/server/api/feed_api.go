package api

import (
	"fmt"
	"sort"
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
	articleTx := services.FilterArticleDraft(database.C)

	if realmId > 0 {
		if realm, err := services.GetRealmWithExtID(uint(realmId)); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("realm was not found: %v", err))
		} else {
			postTx = services.FilterPostWithRealm(postTx, realm.ID)
			articleTx = services.FilterArticleWithRealm(articleTx, realm.ID)
		}
	}

	if len(c.Query("authorId")) > 0 {
		var author models.Account
		if err := database.C.Where(&models.Account{Name: c.Query("authorId")}).First(&author).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		postTx = postTx.Where("author_id = ?", author.ID)
		articleTx = articleTx.Where("author_id = ?", author.ID)
	}

	if len(c.Query("category")) > 0 {
		postTx = services.FilterPostWithCategory(postTx, c.Query("category"))
		articleTx = services.FilterArticleWithCategory(articleTx, c.Query("category"))
	}
	if len(c.Query("tag")) > 0 {
		postTx = services.FilterPostWithTag(postTx, c.Query("tag"))
		articleTx = services.FilterArticleWithTag(articleTx, c.Query("tag"))
	}

	postCountTx := postTx
	articleCountTx := articleTx

	postCount, err := services.CountPost(postCountTx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	articleCount, err := services.CountArticle(articleCountTx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	postItems, err := services.ListPost(postTx, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	articleItems, err := services.ListArticle(articleTx, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var feed []FeedRecord

	encodeToFeed := func(t string, in any, createdAt time.Time) FeedRecord {
		var result map[string]any
		raw, _ := jsoniter.Marshal(in)
		jsoniter.Unmarshal(raw, &result)

		return FeedRecord{
			Type:      t,
			Data:      result,
			CreatedAt: createdAt,
		}
	}

	for _, post := range postItems {
		feed = append(feed, encodeToFeed("post", post, post.CreatedAt))
	}

	for _, article := range articleItems {
		feed = append(feed, encodeToFeed("article", article, article.CreatedAt))
	}

	sort.Slice(feed, func(i, j int) bool {
		return feed[i].CreatedAt.After(feed[j].CreatedAt)
	})

	return c.JSON(fiber.Map{
		"count": postCount + articleCount,
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
	articleTx := services.FilterArticleWithAuthorDraft(database.C, user.ID)

	postCountTx := postTx
	articleCountTx := articleTx

	postCount, err := services.CountPost(postCountTx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	articleCount, err := services.CountArticle(articleCountTx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	postItems, err := services.ListPost(postTx, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	articleItems, err := services.ListArticle(articleTx, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var feed []FeedRecord

	encodeToFeed := func(t string, in any, createdAt time.Time) FeedRecord {
		var result map[string]any
		raw, _ := jsoniter.Marshal(in)
		jsoniter.Unmarshal(raw, &result)

		return FeedRecord{
			Type:      t,
			Data:      result,
			CreatedAt: createdAt,
		}
	}

	for _, post := range postItems {
		feed = append(feed, encodeToFeed("post", post, post.CreatedAt))
	}

	for _, article := range articleItems {
		feed = append(feed, encodeToFeed("article", article, article.CreatedAt))
	}

	sort.Slice(feed, func(i, j int) bool {
		return feed[i].CreatedAt.After(feed[j].CreatedAt)
	})

	return c.JSON(fiber.Map{
		"count": postCount + articleCount,
		"data":  feed,
	})
}
