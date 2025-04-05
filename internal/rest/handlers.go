package rest

import (
	"docker-compose-training/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Handler struct {
	Services *domain.Service
}

func NewHandler(services *domain.Service) *Handler {
	return &Handler{Services: services}
}

func (h *Handler) InitRoutes() *fiber.App {
	router := fiber.New()

	router.Get("/swagger/*", swagger.HandlerDefault)
	router.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Server is up"})
	})

	apiV1 := router.Group("/api/v1")
	{
		storage := apiV1.Group("/storage")
		{
			storage.Post("/", h.PostFile)
			storage.Get("/", h.GetFilesList)
			storage.Get("/:name/content", h.GetFileContent)
			storage.Delete("/:name", h.RemoveFile)

		}
	}

	return router
}
