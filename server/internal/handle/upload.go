package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/utils/upload"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

type Upload struct{}

const maxUploadSize = 10 << 20 // 10 MiB

var allowedImageExtensions = map[string]struct{}{
	".avif": {},
	".gif":  {},
	".jpeg": {},
	".jpg":  {},
	".png":  {},
	".webp": {},
}

func (*Upload) UploadFile(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ReturnError(c, g.ErrFileReceive, err)
		return
	}
	if fileHeader.Size <= 0 || fileHeader.Size > maxUploadSize {
		ReturnError(c, g.ErrFileUpload, "仅支持上传不超过 10 MiB 的图片")
		return
	}
	if _, ok := allowedImageExtensions[strings.ToLower(filepath.Ext(fileHeader.Filename))]; !ok {
		ReturnError(c, g.ErrFileUpload, "仅支持 AVIF、GIF、JPEG、PNG、WebP 图片")
		return
	}

	// 文件存储接口
	oss := upload.NewOSS()
	filepath, _, err := oss.UploadFile(fileHeader)
	if err != nil {
		ReturnError(c, g.ErrFileUpload, err)
		return
	}

	ReturnSuccess(c, filepath)
}
