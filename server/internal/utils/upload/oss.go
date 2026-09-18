// 文件上传服务：当你上传一个文件（如图片）时，文件存储服务（如阿里云 OSS、AWS S3、腾讯云 COS 等）会将文件存储在其服务器上，并返回一个唯一的 URL，指向该文件。
// 把这个url存储在mysql中可以实现后续的文件访问和管理。
package upload

import (
	"errors"
	"io"
	"mime/multipart"
	g "gin-blog/internal/global"
)

// Object Storage Service
// 面向接口编程，OSS可选七牛云，阿里云，腾讯云等，方便切换使用接口  适配器模式
type OSS interface {
	UploadFile(file *multipart.FileHeader) (string,string,error)
	DeleteFile(key string) error
}

func NewOSS() OSS{
	switch g.GetConfig().Upload.OssType{
	case "local" :
		return &Local{}
	case "qiniuyun":
		return &Qiniuyun{}
	case "aliyun":
		return &Aliyun{}
	default:
		return &Local{}
	}
}

// UploadReader uploads an in-memory or archived file. Archive imports do not
// have multipart.FileHeader values for each image, so this keeps that workflow
// inside the same storage provider configured for normal uploads.
func UploadReader(filename string, reader io.Reader) (string, string, error) {
	if uploader, ok := NewOSS().(interface {
		UploadReader(string, io.Reader) (string, string, error)
	}); ok {
		return uploader.UploadReader(filename, reader)
	}

	return "", "", errors.New("当前存储服务不支持导入压缩包中的图片，请将 upload.oss_type 配置为 aliyun")
}
