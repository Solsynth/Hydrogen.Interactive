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

		posts := api.Group("/posts").Name("Posts API")
		{
			posts.Get("/", listPost)
			posts.Get("/:post", getPost)
			posts.Post("/", createPost)
			posts.Post("/:post/react", reactPost)
			posts.Put("/:postId", editPost)
			posts.Delete("/:postId", deletePost)

			posts.Get("/:post/replies", listReplies)
		}

		api.Get("/categories", listCategories)
		api.Post("/categories", newCategory)
		api.Put("/categories/:categoryId", editCategory)
		api.Delete("/categories/:categoryId", deleteCategory)
	}
}
