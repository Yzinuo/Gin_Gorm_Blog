package upload

import (
	"errors"
	"fmt"
	g "gin-blog/internal/global"
	"gin-blog/internal/utils"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"path"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Aliyun struct{}

// UploadReader uploads a file obtained from an archive while preserving the
// same public URL and object-key conventions as UploadFile.
func (*Aliyun) UploadReader(filename string, reader io.Reader) (string, string, error) {
	client, err := oss.New(g.GetConfig().Aliyun.Endpoint, g.GetConfig().Aliyun.AccessKeyID, g.GetConfig().Aliyun.AccessKeySecret)
	if err != nil {
		return "", "", fmt.Errorf("create aliyun OSS client: %w", err)
	}

	bucket, err := client.Bucket(g.GetConfig().Aliyun.Bucket)
	if err != nil {
		return "", "", fmt.Errorf("open aliyun OSS bucket: %w", err)
	}

	filekey := fmt.Sprintf("%s/%d%s%s", g.GetConfig().Aliyun.ImgPath, time.Now().UnixNano(), utils.MD5(filename), path.Ext(filename))
	options := []oss.Option{oss.ObjectACL(oss.ACLPublicRead)}
	if contentType := mime.TypeByExtension(path.Ext(filename)); contentType != "" {
		options = append(options, oss.ContentType(contentType))
	}
	if err = bucket.PutObject(filekey, reader, options...); err != nil {
		return "", "", fmt.Errorf("upload image to aliyun OSS: %w", err)
	}

	publicURL := fmt.Sprintf("https://%s.%s/%s", g.GetConfig().Aliyun.Bucket, g.GetConfig().Aliyun.Endpoint, filekey)
	return publicURL, filekey, nil
}


func (*Aliyun) UploadFile(file *multipart.FileHeader) (string, string, error) {
	client, err := oss.New(g.GetConfig().Aliyun.Endpoint, g.GetConfig().Aliyun.AccessKeyID, g.GetConfig().Aliyun.AccessKeySecret)
	if err != nil {
		log.Printf("function oss.New() Filed, err: %v", err)
		return "", "", errors.New("function oss.New() filed, err: " + err.Error())
	}

	bucket, err := client.Bucket(g.GetConfig().Aliyun.Bucket)
	if err != nil {
		log.Printf("function client.Bucket() Filed, err: %v", err)
		return "", "", errors.New("function client.Bucket() filed, err: " + err.Error())
	}

	f, openErr := file.Open()
	if openErr != nil {
		log.Printf("function file.Open() Filed, err: %v", openErr)
		return "", "", errors.New("function file.Open() filed, err: " + openErr.Error())
	}
	defer f.Close()

	filekey := fmt.Sprintf("%s/%d%s%s", g.GetConfig().Aliyun.ImgPath, time.Now().Unix(), utils.MD5(file.Filename), path.Ext(file.Filename))
	err = bucket.PutObject(filekey, f, oss.ObjectACL(oss.ACLPublicRead))
	if err != nil {
		log.Printf("function bucket.PutObject() Filed, err: %v", err)
		return "", "", errors.New("function bucket.PutObject() filed, err: " + err.Error())
	}
	
	publicURL := fmt.Sprintf("https://%s.%s/%s", g.GetConfig().Aliyun.Bucket, g.GetConfig().Aliyun.Endpoint, filekey)

	return publicURL, filekey, nil
}

// DeleteFile 从阿里云OSS删除文件
func (*Aliyun) DeleteFile(key string) error {
	client, err := oss.New(g.GetConfig().Aliyun.Endpoint, g.GetConfig().Aliyun.AccessKeyID, g.GetConfig().Aliyun.AccessKeySecret)
	if err != nil {
		log.Printf("function oss.New() Filed, err: %v", err)
		return errors.New("function oss.New() filed, err: " + err.Error())
	}

	bucket, err := client.Bucket(g.GetConfig().Aliyun.Bucket)
	if err != nil {
		log.Printf("function client.Bucket() Filed, err: %v", err)
		return errors.New("function client.Bucket() filed, err: " + err.Error())
	}

	err = bucket.DeleteObject(key)
	if err != nil {
		log.Printf("function bucket.DeleteObject() Filed, err: %v", err)
		return errors.New("function bucket.DeleteObject() filed, err: " + err.Error())
	}

	return nil
}
