package api

import (
	"github.com/gofiber/fiber/v2"
)

func MapAPIs(app *fiber.App, baseURL string) {
	api := app.Group(baseURL).Name("API")
	{
		api.Get("/users/me", getUserinfo)
		api.Get("/users/:account", getOthersInfo)
		api.Get("/users/:account/pin", listOthersPinnedPost)

		api.Get("/publishers/:name", getPublisher)

		recommendations := api.Group("/recommendations").Name("Recommendations API")
		{
			recommendations.Get("/", listRecommendationNews)
			recommendations.Get("/friends", listRecommendationFriends)
			recommendations.Get("/shuffle", listRecommendationShuffle)
		}

		stories := api.Group("/stories").Name("Story API")
		{
			stories.Post("/", createStory)
			stories.Put("/:postId", editStory)
		}
		articles := api.Group("/articles").Name("Article API")
		{
			articles.Post("/", createArticle)
			articles.Put("/:postId", editArticle)
		}

		posts := api.Group("/posts").Name("Posts API")
		{
			posts.Get("/", listPost)
			posts.Get("/search", searchPost)
			posts.Get("/minimal", listPostMinimal)
			posts.Get("/drafts", listDraftPost)
			posts.Get("/:postId", getPost)
			posts.Post("/:postId/react", reactPost)
			posts.Post("/:postId/pin", pinPost)
			posts.Delete("/:postId", deletePost)

			posts.Get("/:postId/replies", listPostReplies)
			posts.Get("/:postId/replies/featured", listPostFeaturedReply)
		}

		subscriptions := api.Group("/subscriptions").Name("Subscriptions API")
		{
			subscriptions.Get("/users/:userId", getSubscriptionOnUser)
			subscriptions.Get("/tags/:tagId", getSubscriptionOnTag)
			subscriptions.Get("/categories/:categoryId", getSubscriptionOnCategory)
			subscriptions.Get("/realms/:realmId", getSubscriptionOnRealm)
			subscriptions.Post("/users/:userId", subscribeToUser)
			subscriptions.Post("/tags/:tagId", subscribeToTag)
			subscriptions.Post("/categories/:categoryId", subscribeToCategory)
			subscriptions.Post("/realms/:realmId", subscribeToRealm)
			subscriptions.Delete("/users/:userId", unsubscribeFromUser)
			subscriptions.Delete("/tags/:tagId", unsubscribeFromTag)
			subscriptions.Delete("/categories/:categoryId", unsubscribeFromCategory)
			subscriptions.Delete("/realms/:realmId", unsubscribeFromRealm)
		}

		api.Get("/categories", listCategories)
		api.Get("/categories/:category", getCategory)
		api.Post("/categories", newCategory)
		api.Put("/categories/:categoryId", editCategory)
		api.Delete("/categories/:categoryId", deleteCategory)

		api.Get("/tags", listTags)
		api.Get("/tags/:tag", getTag)

		api.Get("/whats-new", getWhatsNew)
	}
}
