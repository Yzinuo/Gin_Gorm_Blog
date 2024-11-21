package upload

import (
	"errors"
	"fmt"
	g "gin-blog/internal/global"
	"gin-blog/internal/utils"
	"log"
	"mime/multipart"
	"path"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Aliyun struct{}


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