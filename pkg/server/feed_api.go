package server

import "C"
import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type FeedItem struct {
	models.BaseModel

	Title         string `json:"title"`
	Description   string `json:"description"`
	Content       string `json:"content"`
	ModelType     string `json:"model_type"`
	CommentCount  int64  `json:"comment_count"`
	ReactionCount int64  `json:"reaction_count"`
	AuthorID      uint   `json:"author_id"`
	RealmID       *uint  `json:"realm_id"`

	Author models.Account `json:"author" gorm:"embedded"`
}

const (
	queryArticle = "id, created_at, updated_at, title, content, description, realm_id, author_id, 'article' as model_type"
	queryMoment  = "id, created_at, updated_at, NULL as title, content, NULL as description, realm_id, author_id, 'moment' as model_type"
)

func listFeed(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)

	if take > 20 {
		take = 20
	}

	var whereCondition string

	if realmId > 0 {
		whereCondition += fmt.Sprintf("feed.realm_id = %d", realmId)
	} else {
		whereCondition += "feed.realm_id IS NULL"
	}

	var author models.Account
	if len(c.Query("authorId")) > 0 {
		if err := database.C.Where(&models.Account{Name: c.Query("authorId")}).First(&author).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		} else {
			whereCondition += fmt.Sprintf("AND feed.author_id = %d", author.ID)
		}
	}

	var result []FeedItem

	userTable := viper.GetString("database.prefix") + "accounts"
	commentTable := viper.GetString("database.prefix") + "comments"
	reactionTable := viper.GetString("database.prefix") + "reactions"

	database.C.Raw(fmt.Sprintf(`SELECT feed.*, author.*, 
		COALESCE(comment_count, 0) as comment_count, 
		COALESCE(reaction_count, 0) as reaction_count
		FROM (? UNION ALL ?) as feed
		INNER JOIN %s as author ON author_id = author.id
		LEFT JOIN (SELECT article_id, moment_id, COUNT(*) as comment_count
            FROM %s
            GROUP BY article_id, moment_id) as comments
            ON (feed.model_type = 'article' AND feed.id = comments.article_id) OR 
			   (feed.model_type = 'moment' AND feed.id = comments.moment_id)
        LEFT JOIN (SELECT article_id, moment_id, COUNT(*) as reaction_count
        	FROM %s
            GROUP BY article_id, moment_id) as reactions
            ON (feed.model_type = 'article' AND feed.id = reactions.article_id) OR 
			   (feed.model_type = 'moment' AND feed.id = reactions.moment_id)
		WHERE %s LIMIT ? OFFSET ?`, userTable, commentTable, reactionTable, whereCondition),
		database.C.Select(queryArticle).Model(&models.Article{}),
		database.C.Select(queryMoment).Model(&models.Moment{}),
		take,
		offset,
	).Scan(&result)

	return c.JSON(result)
}
