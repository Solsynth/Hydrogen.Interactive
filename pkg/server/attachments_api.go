package server

import (
	"path/filepath"

	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func readAttachment(c *fiber.Ctx) error {
	id := c.Params("fileId")
	basepath := viper.GetString("content")

	return c.SendFile(filepath.Join(basepath, id))
}

func uploadAttachment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	file, err := c.FormFile("attachment")
	if err != nil {
		return err
	}

	attachment, err := services.NewAttachment(user, file)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.SaveFile(file, attachment.GetStoragePath()); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"info": attachment,
		"url":  attachment.GetAccessPath(),
	})
}

func deleteAttachment(c *fiber.Ctx) error {
	id := c.Params("fileId")
	user := c.Locals("principal").(models.Account)

	attachment, err := services.GetAttachmentByUUID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if attachment.AuthorID != user.ID {
		return fiber.NewError(fiber.StatusNotFound, "record not created by you")
	}

	if err := services.DeleteAttachment(attachment); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}
