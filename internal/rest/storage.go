package rest

import (
	"docker-compose-training/internal/domain"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"os"
)

// PostFile uploads file to storage
//
//	@Summary      Uploads file
//	@Description  uploads file
//	@Tags         files
//	@Accept       json
//	@Produce      json
//
// @Param			input	formData	file	true	"Files"
// @Param			name  query       string             false    "name"
// @Success      200
// @Failure      400  {string}  string
// @Failure      500  {string}  string
// @Router       /storage [post]
func (h *Handler) PostFile(c *fiber.Ctx) error {
	ctx := c.Context()
	file, err := c.FormFile("input")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Errorf("failed to get file: %w", err).Error())
	}
	slog.Info("got file", "filename", file.Filename)

	name := c.Query("name")
	if name == "" {
		name = file.Filename
	}

	filePath := fmt.Sprintf("%s%s", h.basePath, name)
	if _, err = os.CreateTemp(h.basePath, name); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Errorf("failed to create tmp file %s: %w", filePath, err).Error())
	}

	err = c.SaveFile(file, filePath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Errorf("failed to save file to tmp folder: %w", err).Error())
	}
	defer func() {
		err = os.Remove(filePath)
		if err != nil {
			slog.Info("failed to remove temporary file", "path", h.basePath)
		}
	}()

	if err = h.Services.PostFile(ctx, name); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

// GetFilesList get names of all uploaded files
//
//	@Summary      Get names
//	@Description  get names of all uploaded files
//	@Tags         files
//	@Accept       json
//	@Produce      json
//
// @Success      200  {object}  FileNamesResponse
// @Failure      500  {string}  string
// @Router       /storage [get]
func (h *Handler) GetFilesList(c *fiber.Ctx) error {
	ctx := c.Context()
	files, err := h.Services.GetFilesList(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(FileNamesResponse{Names: files})
}

// GetFileContent get text content of the file
//
//	@Summary      Get content
//	@Description  get text content of the file
//	@Tags         files
//	@Accept       json
//	@Produce      json
//
// @Param		name  path       string             true    "name"
// @Success      200  {string}  string
// @Failure      500  {string}  string
// @Router       /storage/{name}/content [get]
func (h *Handler) GetFileContent(c *fiber.Ctx) error {
	ctx := c.Context()
	fileName := c.Params("name")

	content, err := h.Services.GetFileContent(ctx, fileName)
	if errors.Is(err, domain.NoObjectErr) {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString(content.Content)
}

// RemoveFile remove existing file from storage
//
//	@Summary      Remove file
//	@Description  remove existing file from storage
//	@Tags         files
//	@Accept       json
//	@Produce      json
//
// @Param		name  path       string             true    "name"
// @Success      200
// @Failure      500  {string}  string
// @Router       /storage/{name} [delete]
func (h *Handler) RemoveFile(c *fiber.Ctx) error {
	ctx := c.Context()
	fileName := c.Params("name")

	err := h.Services.RemoveFile(ctx, fileName)
	if errors.Is(err, domain.NoObjectErr) {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
