package api

import (
	"github.com/gofiber/fiber/v2"
)

func MapAPIs(app *fiber.App) {
	api := app.Group("/api").Name("API")
	{
		api.Get("/users/me", getUserinfo)
		api.Get("/users/:accountId", getOthersInfo)

		api.Get("/feed", listFeed)

		drafts := api.Group("/drafts").Name("Draft box API")
		{
			drafts.Get("/posts", listDraftPost)
			drafts.Get("/articles", listDraftArticle)
		}

		posts := api.Group("/posts").Name("Posts API")
		{
			posts.Get("/", listPost)
			posts.Get("/:post", getPost)
			posts.Post("/", createPost)
			posts.Post("/:post/react", reactPost)
			posts.Put("/:postId", editPost)
			posts.Delete("/:postId", deletePost)

			posts.Get("/:post/replies", listPostReplies)
		}

		articles := api.Group("/articles").Name("Articles API")
		{
			articles.Get("/", listArticle)
			articles.Get("/:article", getArticle)
			articles.Post("/", createArticle)
			articles.Post("/:article/react", reactArticle)
			articles.Put("/:articleId", editArticle)
			articles.Delete("/:articleId", deleteArticle)
		}

		api.Get("/categories", listCategories)
		api.Post("/categories", newCategory)
		api.Put("/categories/:categoryId", editCategory)
		api.Delete("/categories/:categoryId", deleteCategory)
	}
}
