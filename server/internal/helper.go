// helper的目的就是定义加载config的函数
package ginblog

import (
	"context"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitLogger(conf *g.Config) {
	var level slog.Level
	switch conf.Log.Level{
	case "debug" :
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}

	option := &slog.HandlerOptions{
		AddSource: false, //是否显示源文件的信息
		Level : level,
		ReplaceAttr: func(groups []string, a slog.Attr)slog.Attr{
	            if a.Key == slog.TimeKey{
					if t,ok := a.Value.Any().(time.Time); ok {
						a.Value = slog.StringValue(t.Format(time.DateTime))
					}
				}
				return a
		},
	}

	var handler slog.Handler
	switch conf.Log.Format{
	case "json":
		handler = slog.NewJSONHandler(os.Stdout,option)
	case "text":
		handler = slog.NewTextHandler(os.Stdout,option)
	default:
		handler = slog.NewTextHandler(os.Stdout,option)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger) // 设定默认的slog
}

func InitDatabase(conf *g.Config) *gorm.DB{
	dbType := conf.DbType()
	dsn := conf.DbDSN()

	var level logger.LogLevel
	switch conf.Log.Level{
	case "silent":
		level = logger.Silent
	case "error":
		level = logger.Error
	case "warn":
		level = logger.Warn
	case "info":
		level = logger.Info
	default:
		level = logger.Error
	}

	config := &gorm.Config{
		Logger: logger.Default.LogMode(level),
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
		SkipDefaultTransaction:                   true, // 禁用默认事务（提高运行速度）
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 单数表名
		},
	}

	var db *gorm.DB
	var err error
	switch dbType {
	case "mysql":
		db,err = gorm.Open(mysql.Open(dsn),config)
	case "sqlite":
		db,err = gorm.Open(sqlite.Open(dsn),config)
	default:
		log.Fatalf("不支持的数据库类型")
	}
	
	if err != nil{
		log.Fatal("连接数据库失败")
	}
	slog.Info("连接数据库成功")

	if conf.Server.DbAutoMigrate{
		if err := model.MakeMigrate(db);err != nil{
			log.Fatalf("数据库迁移失败",err)
		}
		slog.Info("数据库自动迁移成功")
	}
	return db
}

func InitRedis(conf *g.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB: conf.Redis.DB,
	})

    _,err := rdb.Ping(context.Background()).Result()
	if err != nil{
		log.Fatal("连接redis失败")
	}

	slog.Info("连接redis成功")
	return rdb
}