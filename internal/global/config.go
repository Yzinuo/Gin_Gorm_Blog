package global

type Config struct {
	JWT struct {
		Secret string
		Expire int64 // hour
		Issuer string
	}

	Upload  struct{
		OssType  string // local | qniuyun
		Path  	 string // 本地访问路径
		StorePath  string  // 本地存储路径
	}

	Qiniu struct {
		ImgPath       string // 外链链接
		Zone          string // 存储区域
		Bucket        string // 空间名称
		AccessKey     string // 秘钥AK
		SecretKey     string // 秘钥SK
		UseHTTPS      bool   // 是否使用https
		UseCdnDomains bool   // 上传是否使用 CDN 上传加速
	}
}

var Conf *Config