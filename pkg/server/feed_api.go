package server

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
	"strings"

	"git.solsynth.dev/hydrogen/interactive/pkg/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

const (
	queryArticle = "id, created_at, updated_at, alias, title, NULL as content, description, realm_id, author_id, 'article' as model_type"
	queryMoment  = "id, created_at, updated_at, alias, NULL as title, content, NULL as description, realm_id, author_id, 'moment' as model_type"
)

func listFeed(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmAlias := c.Query("realm")

	if take > 20 {
		take = 20
	}

	var whereConditions []string

	if len(realmAlias) > 0 {
		realm, err := services.GetRealmWithAlias(realmAlias)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("related realm was not found: %v", err))
		}
		whereConditions = append(whereConditions, fmt.Sprintf("feed.realm_id = %d", realm.ID))
	}

	var author models.Account
	if len(c.Query("authorId")) > 0 {
		if err := database.C.Where(&models.Account{Name: c.Query("authorId")}).First(&author).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		} else {
			whereConditions = append(whereConditions, fmt.Sprintf("feed.author_id = %d", author.ID))
		}
	}

	var whereStatement string
	if len(whereConditions) > 0 {
		whereStatement += "WHERE " + strings.Join(whereConditions, " AND ")
	}

	var result []*models.Feed

	userTable := viper.GetString("database.prefix") + "accounts"
	commentTable := viper.GetString("database.prefix") + "comments"
	reactionTable := viper.GetString("database.prefix") + "reactions"

	database.C.Raw(
		fmt.Sprintf(`SELECT feed.*, author.*,
		COALESCE(comment_count, 0) AS comment_count, 
		COALESCE(reaction_count, 0) AS reaction_count
		FROM (? UNION ALL ?) AS feed
		INNER JOIN %s AS author ON author_id = author.id
		LEFT JOIN (SELECT article_id, moment_id, COUNT(*) AS comment_count
            FROM %s
            GROUP BY article_id, moment_id) AS comments
            ON (feed.model_type = 'article' AND feed.id = comments.article_id) OR 
			   (feed.model_type = 'moment' AND feed.id = comments.moment_id)
        LEFT JOIN (SELECT article_id, moment_id, COUNT(*) AS reaction_count
        	FROM %s
            GROUP BY article_id, moment_id) AS reactions
            ON (feed.model_type = 'article' AND feed.id = reactions.article_id) OR 
			   (feed.model_type = 'moment' AND feed.id = reactions.moment_id)
		%s ORDER BY feed.created_at desc  LIMIT ? OFFSET ?`,
			userTable,
			commentTable,
			reactionTable,
			whereStatement,
		),
		database.C.Select(queryArticle).Model(&models.Article{}),
		database.C.Select(queryMoment).Model(&models.Moment{}),
		take,
		offset,
	).Scan(&result)

	if !c.QueryBool("noReact", false) {
		var reactions []struct {
			PostID uint
			Symbol string
			Count  int64
		}

		revertReaction := func(dataset string) error {
			itemMap := lo.SliceToMap(lo.FilterMap(result, func(item *models.Feed, index int) (*models.Feed, bool) {
				return item, item.ModelType == dataset
			}), func(item *models.Feed) (uint, *models.Feed) {
				return item.ID, item
			})

			idx := lo.Map(lo.Filter(result, func(item *models.Feed, index int) bool {
				return item.ModelType == dataset
			}), func(item *models.Feed, index int) uint {
				return item.ID
			})

			if err := database.C.Model(&models.Reaction{}).
				Select(dataset+"_id as post_id, symbol, COUNT(id) as count").
				Where(dataset+"_id IN (?)", idx).
				Group("post_id, symbol").
				Scan(&reactions).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			list := map[uint]map[string]int64{}
			for _, info := range reactions {
				if _, ok := list[info.PostID]; !ok {
					list[info.PostID] = make(map[string]int64)
				}
				list[info.PostID][info.Symbol] = info.Count
			}

			for k, v := range list {
				if post, ok := itemMap[k]; ok {
					post.ReactionList = v
				}
			}

			return nil
		}

		if err := revertReaction("article"); err != nil {
			return err
		}
		if err := revertReaction("moment"); err != nil {
			return err
		}
	}

	if !c.QueryBool("noAttachment", false) {
		revertAttachment := func(dataset string) error {
			var attachments []struct {
				models.Attachment

				PostID uint `json:"post_id"`
			}

			itemMap := lo.SliceToMap(lo.FilterMap(result, func(item *models.Feed, index int) (*models.Feed, bool) {
				return item, item.ModelType == dataset
			}), func(item *models.Feed) (uint, *models.Feed) {
				return item.ID, item
			})

			idx := lo.Map(lo.Filter(result, func(item *models.Feed, index int) bool {
				return item.ModelType == dataset
			}), func(item *models.Feed, index int) uint {
				return item.ID
			})

			if err := database.C.
				Model(&models.Attachment{}).
				Select(dataset+"_id as post_id, *").
				Where(dataset+"_id IN (?)", idx).
				Scan(&attachments).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			list := map[uint][]models.Attachment{}
			for _, info := range attachments {
				list[info.PostID] = append(list[info.PostID], info.Attachment)
			}

			for k, v := range list {
				if post, ok := itemMap[k]; ok {
					post.Attachments = v
				}
			}

			return nil
		}

		if err := revertAttachment("article"); err != nil {
			return err
		}
		if err := revertAttachment("moment"); err != nil {
			return err
		}
	}

	var count int64
	database.C.Raw(`SELECT COUNT(*) FROM (? UNION ALL ?) as feed`,
		database.C.Select(queryArticle).Model(&models.Article{}),
		database.C.Select(queryMoment).Model(&models.Moment{}),
	).Scan(&count)

	return c.JSON(fiber.Map{
		"count": count,
		"data":  result,
	})
}
