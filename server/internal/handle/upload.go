package handle

import (
	"fmt"
	g "gin-blog/internal/global"
	"gin-blog/internal/utils/upload"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"
)

type Upload struct{}

const maxUploadSize = 10 << 20 // 10 MiB

var genericImageMIMEByExt = map[string]string{
	".avif": "image/avif",
	".gif":  "image/gif",
	".jpeg": "image/jpeg",
	".jpg":  "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
}

func (*Upload) UploadFile(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize+(1<<20))
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		ReturnError(c, g.ErrFileReceive, err)
		return
	}
	purpose := strings.ToLower(c.PostForm("purpose"))
	if purpose == "" {
		purpose = "article"
	}
	maxWidth := 1600
	switch purpose {
	case "article":
	case "avatar":
		maxWidth = 512
	case "cover":
		maxWidth = 1920
	default:
		ReturnError(c, g.ErrRequest, "purpose 必须为 article、avatar 或 cover")
		return
	}

	file, err := spoolUpload(header, maxUploadSize, genericImageMIMEByExt)
	if err != nil {
		ReturnError(c, g.ErrFileUpload, err)
		return
	}
	defer os.Remove(file.path)
	store, err := upload.NewObjectStore()
	if err != nil {
		ReturnError(c, g.ErrFileUpload, err)
		return
	}
	prefix := fmt.Sprintf("uploads/%s/%s/%s/%s", purpose, time.Now().Format("2006/01"), file.hash[:2], file.hash)
	sourceKey := prefix + "/source" + file.ext
	source, err := uploadPrepared(c, store, file, sourceKey, purpose+".source")
	if err != nil {
		ReturnError(c, g.ErrFileUpload, err)
		return
	}
	publicURL := source.URL

	// GIF animation and AVIF inputs are retained without lossy transcoding.
	if file.mime != "image/gif" && file.mime != "image/avif" {
		img, _, decodeErr := decodeImageFile(file)
		if decodeErr != nil {
			_ = store.Delete(c, sourceKey)
			ReturnError(c, g.ErrFileUpload, decodeErr)
			return
		}
		variant, encodeErr := encodeWebPVariant(img, maxWidth, 82)
		if encodeErr != nil {
			_ = store.Delete(c, sourceKey)
			ReturnError(c, g.ErrFileUpload, encodeErr)
			return
		}
		defer os.Remove(variant.path)
		deliveryKey := fmt.Sprintf("%s/delivery-v1-%d.webp", prefix, variant.width)
		delivery, uploadErr := uploadPrepared(c, store, variant, deliveryKey, purpose+".delivery")
		if uploadErr != nil {
			_ = store.Delete(c, sourceKey)
			ReturnError(c, g.ErrFileUpload, uploadErr)
			return
		}
		publicURL = delivery.URL
	}

	ReturnSuccess(c, publicURL)
}
