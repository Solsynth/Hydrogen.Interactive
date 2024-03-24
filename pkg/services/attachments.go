package services

import (
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"git.solsynth.dev/hydrogen/interactive/pkg/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func GetAttachmentByUUID(fileId string) (models.Attachment, error) {
	var attachment models.Attachment
	if err := database.C.Where(models.Attachment{
		FileID: fileId,
	}).First(&attachment).Error; err != nil {
		return attachment, err
	}
	return attachment, nil
}

func NewAttachment(user models.Account, header *multipart.FileHeader) (models.Attachment, error) {
	attachment := models.Attachment{
		FileID:   uuid.NewString(),
		Filesize: header.Size,
		Filename: header.Filename,
		Mimetype: "unknown/unknown",
		Type:     models.AttachmentOthers,
		AuthorID: user.ID,
	}

	// Open file
	file, err := header.Open()
	if err != nil {
		return attachment, err
	}
	defer file.Close()

	// Detect mimetype
	fileHeader := make([]byte, 512)
	_, err = file.Read(fileHeader)
	if err != nil {
		return attachment, err
	}
	attachment.Mimetype = http.DetectContentType(fileHeader)

	switch strings.Split(attachment.Mimetype, "/")[0] {
	case "image":
		attachment.Type = models.AttachmentPhoto
	case "video":
		attachment.Type = models.AttachmentVideo
	case "audio":
		attachment.Type = models.AttachmentAudio
	default:
		attachment.Type = models.AttachmentOthers
	}

	// Save into database
	err = database.C.Save(&attachment).Error

	return attachment, err
}

func DeleteAttachment(item models.Attachment) error {
	if err := database.C.Delete(&item).Error; err != nil {
		return err
	} else {
		basepath := viper.GetString("content")
		fullpath := filepath.Join(basepath, item.FileID)

		os.Remove(fullpath)
	}

	return nil
}
