package models

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/samber/lo"
)

const (
	MaxMangaPreviewFileSize = (1 << 20) * 2  // 2 MB
	MaxMangaPageZipFileSize = (1 << 20) * 20 // 20 MB
)

var (
	ErrInvalidFileExt = errors.New("invalid file extension")
)

type UploadFile struct {
	Data     []byte
	Filename string
}

// ValidateExtension checks if file extension is valid.
// Extensions format example: .png .jpg .gif
func ValidateExtension(filename string, extensions ...string) error {
	if !lo.Contains[string](extensions, filepath.Ext(filename)) {
		return fmt.Errorf("Allowed file extensions: %s", extensions)
	}

	return nil
}
