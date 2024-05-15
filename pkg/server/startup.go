package server

import (
	"git.solsynth.dev/hydrogen/interactive/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

var A *fiber.App

func NewServer() {
	templates := html.NewFileSystem(http.FS(pkg.FS), ".gohtml")

	A = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		ServerHeader:          "Hydrogen.Interactive",
		AppName:               "Hydrogen.Interactive",
		ProxyHeader:           fiber.HeaderXForwardedFor,
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		BodyLimit:             50 * 1024 * 1024,
		EnablePrintRoutes:     viper.GetBool("debug.print_routes"),
		Views:                 templates,
		ViewsLayout:           "views/index",
	})

	A.Use(idempotency.New())
	A.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodOptions,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	A.Use(logger.New(logger.Config{
		Format: "${status} | ${latency} | ${method} ${path}\n",
		Output: log.Logger,
	}))

	A.Get("/.well-known", getMetadata)

	api := A.Group("/api").Name("API")
	{
		api.Get("/users/me", authMiddleware, getUserinfo)
		api.Get("/users/:accountId", getOthersInfo)

		api.Get("/attachments/o/:fileId", readAttachment)
		api.Post("/attachments", authMiddleware, uploadAttachment)
		api.Delete("/attachments/:id", authMiddleware, deleteAttachment)

		api.Get("/feed", listFeed)

		posts := api.Group("/posts").Name("Posts API")
		{
			posts.Get("/", listPost)
			posts.Get("/:postId", getPost)
			posts.Post("/", authMiddleware, createPost)
			posts.Post("/:postId/react", authMiddleware, reactPost)
			posts.Put("/:postId", authMiddleware, editPost)
			posts.Delete("/:postId", authMiddleware, deletePost)

			posts.Get("/:postId/replies", listReplies)
		}

		api.Get("/categories", listCategories)
		api.Post("/categories", authMiddleware, newCategory)
		api.Put("/categories/:categoryId", authMiddleware, editCategory)
		api.Delete("/categories/:categoryId", authMiddleware, deleteCategory)
	}

	A.Use(favicon.New(favicon.Config{
		FileSystem: http.FS(pkg.FS),
		File:       "views/favicon.png",
		URL:        "/favicon.png",
	}))

	A.Get("/", func(c *fiber.Ctx) error {
		return c.Render("views/open", fiber.Map{
			"frontend": viper.GetString("frontend"),
		})
	})
}

func Listen() {
	if err := A.Listen(viper.GetString("bind")); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when starting server...")
	}
}
