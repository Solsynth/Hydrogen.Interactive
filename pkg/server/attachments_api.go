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
	hashcode := c.FormValue("hashcode")
	if len(hashcode) != 64 {
		return fiber.NewError(fiber.StatusBadRequest, "please provide a SHA256 hashcode, length should be 64 characters")
	}
	file, err := c.FormFile("attachment")
	if err != nil {
		return err
	}

	attachment, err := services.NewAttachment(user, file, hashcode)
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
	id, _ := c.ParamsInt("id", 0)
	user := c.Locals("principal").(models.Account)

	attachment, err := services.GetAttachmentByID(uint(id))
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
