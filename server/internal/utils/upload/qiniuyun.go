package upload

import (
	"context"
	"errors"
	"fmt"
	g "gin-blog/internal/global"
	"gin-blog/internal/utils"
	"log/slog"
	"mime/multipart"
	"path"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Qiniuyun struct{}


func (*Qiniuyun) UploadFile(file *multipart.FileHeader) (string,string,error){
	putPolicy := storage.PutPolicy{Scope:g.GetConfig().Qiniu.Bucket} // 指定数据存储的位值
	mac := qbox.NewMac(g.GetConfig().Qiniu.AccessKey,g.GetConfig().Qiniu.SecretKey)
	uptoken := putPolicy.UploadToken(mac) // qniuyun sdk 上传需要token
	formUploader := storage.NewFormUploader(QiniuyunConfig()) // 文件上传的配置
	ret := storage.PutRet{}// 上传请求的回复
	putExtra := storage.PutExtra{Params: map[string]string{"x:name": "gvb image"}}//上传额外参数 七牛云存储都会在文件名前加上x:name，方便区分

	f,openErr := file.Open()
	if openErr != nil {
		slog.Error("function file.Open() Filed", slog.Any("err",openErr.Error()))
		return "","",errors.New("function file.Open() filed,err:" + openErr.Error())
	}
	defer f.Close()

	filekey := fmt.Sprintf("%d%s%s",time.Now().Unix(),utils.MD5(file.Filename),path.Ext(file.Filename)) // 确保文件名格式的唯一性
	puterr := formUploader.Put(context.Background(),&ret,uptoken,filekey,f,file.Size,&putExtra)
	if puterr != nil {
		slog.Error("function formUploader.Put() Filed", slog.Any("err",puterr.Error()))
		return "","",errors.New("function formUploader.Put() filed,err:" + puterr.Error())
	}

	return g.GetConfig().Qiniu.ImgPath + "/" + ret.Key,ret.Key,nil
}

func (*Qiniuyun) DeleteFile(key string) error{
	mac := qbox.NewMac(g.GetConfig().Qiniu.AccessKey,g.GetConfig().Qiniu.SecretKey)
	cfg := QiniuyunConfig()
	bucketManager := storage.NewBucketManager(mac,cfg)
	if err := bucketManager.Delete(g.GetConfig().Qiniu.Bucket,key); err!= nil {
		slog.Error("function bucketManager.Delete() Filed", slog.Any("err",err.Error()))
		return errors.New("function bucketManager.Delete() filed,err:" + err.Error())
	}
	return nil
}

func QiniuyunConfig() *storage.Config{
	cfg := storage.Config{}

	cfg.UseHTTPS = g.GetConfig().Qiniu.UseHTTPS
	cfg.UseCdnDomains = g.GetConfig().Qiniu.UseCdnDomains

	switch g.GetConfig().Qiniu.Zone{
	case "ZoneHuadong":
		cfg.Zone = &storage.ZoneHuabei
	case "ZoneHuabei":
		cfg.Zone = &storage.ZoneHuabei
	case "ZoneHuanan":
		cfg.Zone = &storage.ZoneHuanan
	case "ZoneBeimei":
		cfg.Zone = &storage.ZoneBeimei
	case "ZoneXinjiapo":
		cfg.Zone = &storage.ZoneXinjiapo
	}

	return &cfg
}