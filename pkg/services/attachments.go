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

func GetAttachmentByID(id uint) (models.Attachment, error) {
	var attachment models.Attachment
	if err := database.C.Where(models.Attachment{
		BaseModel: models.BaseModel{ID: id},
	}).First(&attachment).Error; err != nil {
		return attachment, err
	}
	return attachment, nil
}

func GetAttachmentByUUID(fileId string) (models.Attachment, error) {
	var attachment models.Attachment
	if err := database.C.Where(models.Attachment{
		FileID: fileId,
	}).First(&attachment).Error; err != nil {
		return attachment, err
	}
	return attachment, nil
}

func GetAttachmentByHashcode(hashcode string) (models.Attachment, error) {
	var attachment models.Attachment
	if err := database.C.Where(models.Attachment{
		Hashcode: hashcode,
	}).First(&attachment).Error; err != nil {
		return attachment, err
	}
	return attachment, nil
}

func NewAttachment(user models.Account, header *multipart.FileHeader, hashcode string) (models.Attachment, error) {
	var attachment models.Attachment
	existsAttachment, err := GetAttachmentByHashcode(hashcode)
	if err != nil {
		// Upload the new file
		attachment = models.Attachment{
			FileID:   uuid.NewString(),
			Filesize: header.Size,
			Filename: header.Filename,
			Hashcode: hashcode,
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
	} else {
		// Instant upload, build link with the exists file
		attachment = models.Attachment{
			FileID:   existsAttachment.FileID,
			Filesize: header.Size,
			Filename: header.Filename,
			Hashcode: hashcode,
			Mimetype: existsAttachment.Mimetype,
			Type:     existsAttachment.Type,
			AuthorID: user.ID,
		}
	}

	// Save into database
	err = database.C.Save(&attachment).Error

	return attachment, err
}

func DeleteAttachment(item models.Attachment) error {
	var dupeCount int64
	if err := database.C.
		Where(&models.Attachment{Hashcode: item.Hashcode}).
		Model(&models.Attachment{}).
		Count(&dupeCount).Error; err != nil {
		dupeCount = -1
	}

	if err := database.C.Delete(&item).Error; err != nil {
		return err
	}

	if dupeCount != -1 && dupeCount <= 1 {
		// Safe for deletion the physics file
		basepath := viper.GetString("content")
		fullpath := filepath.Join(basepath, item.FileID)

		os.Remove(fullpath)
	}

	return nil
}
