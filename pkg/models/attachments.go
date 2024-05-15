package models

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type AttachmentType = uint8

const (
	AttachmentOthers = AttachmentType(iota)
	AttachmentPhoto
	AttachmentVideo
	AttachmentAudio
)

type Attachment struct {
	BaseModel

	FileID   string         `json:"file_id"`
	Filesize int64          `json:"filesize"`
	Filename string         `json:"filename"`
	Mimetype string         `json:"mimetype"`
	Hashcode string         `json:"hashcode"`
	Type     AttachmentType `json:"type"`
	Author   Account        `json:"author"`
	AuthorID uint           `json:"author_id"`

	PostID *uint `json:"post_id"`
}

func (v Attachment) GetStoragePath() string {
	basepath := viper.GetString("content")
	return filepath.Join(basepath, v.FileID)
}

func (v Attachment) GetAccessPath() string {
	return fmt.Sprintf("/api/attachments/o/%s", v.FileID)
}
