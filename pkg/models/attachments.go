package models

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
)

type Attachment struct {
	BaseModel

	FileID   string  `json:"file_id"`
	Filesize int64   `json:"filesize"`
	Filename string  `json:"filename"`
	Mimetype string  `json:"mimetype"`
	Post     *Post   `json:"post"`
	Author   Account `json:"author"`
	PostID   *uint   `json:"post_id"`
	AuthorID uint    `json:"author_id"`
}

func (v Attachment) GetStoragePath() string {
	basepath := viper.GetString("content")
	return filepath.Join(basepath, v.FileID)
}

func (v Attachment) GetAccessPath() string {
	return fmt.Sprintf("/api/attachments/o/%s", v.FileID)
}
