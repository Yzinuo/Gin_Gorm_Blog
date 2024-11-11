package main

import (
	"flag"
	"log"
	"strings"

	ginblog "gin-blog/internal"
	g "gin-blog/internal/global"
	"gin-blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main(){
	// 根据命令行给出的文件路径，读取配置文件
	configPath := flag.String("c","../config.yml","配置文件路径")
	flag.Parse()
	conf := g.ReadConfig(*configPath)

	//初始化gin
	ginblog.InitLogger(conf)
	db := ginblog.InitDatabase(conf)
	rdb := ginblog.InitRedis(conf)
	gin.SetMode(conf.Server.Mode) // 不同的mode会影响日志输出
	r := gin.New()
	// 定义信任的代理服务器
	r.SetTrustedProxies([]string{"*"})
	r.Use(gin.Logger(),gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.WithCookieStore(conf.Session.Name,conf.Session.Salt))
	r.Use(middleware.WithGormDB(db))
	r.Use(middleware.WithRDB(rdb))
	ginblog.RegisterAllHandler(r) // 打包注册所有的handler

	// 本地存储，需要设置静态文件服务
	if conf.Upload.OssType == "local" {
		r.Static(conf.Upload.Path,conf.Upload.StorePath)
	}
	
	// 定义程序运行的ip地址和端口
	serverAddr := conf.Server.Port
	if serverAddr[0] == ':' || strings.HasPrefix(serverAddr,"0.0.0.0:") { // 没有配域名
			log.Printf("Serving HTTP on (http://localhost:%s/) ... \n",strings.Split(serverAddr, ":")[1])
	}else{ // 配了域名
			log.Printf("Serving HTTP on (http//%s/)\n",serverAddr)
	}
	
	r.Run(serverAddr)
}