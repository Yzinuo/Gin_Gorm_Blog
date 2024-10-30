package upload

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/utils"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

type Local struct{}

// 上传文件到本地  对于本机来说 存储路径的作用就是类似url
func (*Local) UploadFile(file *multipart.FileHeader) (string,string,error){
	ext := path.Ext(file.Filename) // 获取文件后缀 （.jpg  .png）
	name := strings.TrimSuffix(file.Filename,ext) // 去掉后缀，得到文件名
	name = utils.MD5(name) // HASH 文件名
    filename := name + "_" + time.Now().Format("20060102150405") + ext // 文件名 + 时间戳  新名字

	conf := g.Conf.Upload

	MKdir := os.MkdirAll(conf.StorePath,os.ModePerm) // 创建文件夹
	if MKdir != nil {
		slog.Error("function os.MkdirAll() Filed", slog.Any("err",MKdir.Error()))
		return "","",errors.New("function os.MkdirAll() Filed,err :"+MKdir.Error())
	}
	// 文件存储路径
	storePath := conf.StorePath + "/" + filename
	filePath := conf.Path + "/" + filename
	
	f,err  := file.Open()
	if err != nil {
		slog.Error("function filepath.Open() Filed", slog.Any("err",err.Error()))
		return "","",errors.New("function filepath.Open() filed ,err:" + err.Error())
	}
	defer  f.Close()

	out,err := os.Create(storePath)
	if err != nil {
		slog.Error("function os.Create() Filed", slog.Any("err",err.Error()))
		return "","",errors.New("function os.Create() filed,err:" + err.Error())
	}
	defer out.Close()

	_,err = io.Copy(out,f)
	if err != nil {
		slog.Error("function io.copy() Filed", slog.Any("err",err.Error()))
		return "","",errors.New("function io.copy() filed ,err:" + err.Error())
	}

	return filePath,filename,nil
}

func (*Local) DeleteFile(key string) error{
	p := g.GetConfig().Upload.StorePath + "/" + key
	// 防止有人恶心更改参数：key 中包含 ../ 等路径遍历符号来访问或写入其他目录的文件。
	// 确保是正确的删除路径
	if strings.Contains(p,g.GetConfig().Upload.StorePath){
		if err := os.Remove(p); err != nil{
			return errors.New("本地文件删除失败,err:"+err.Error())
		}
	}
	return nil
}