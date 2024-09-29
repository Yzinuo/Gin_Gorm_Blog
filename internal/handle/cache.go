// 缓存由redis 内存数据库来管理
// 存储在缓存中的数据： 1. 频繁访问的数据  2.静态数据：如配置等等
package handle

import (
	"context"
	"encoding/json"
	"errors"
	model "gin-blog/internal/model"
	g "gin-blog/internal/global"

	"github.com/go-redis/redis/v9"
)

var rdbctx  = context.Background()

// 页面背景
func AddPageCache(rdb *redis.Client,page []model.Page) error{
	data,err := json.Marshal(page)
	if err != nil {
		return errors.New("JSON fail to murshal")	
	}
	
	return rdb.Set(rdbctx,g.PAGE,string(data),0).Err()
}

func DelPageCache(rdb *redis.Client) error{
	return rdb.Del(rdbctx,g.PAGE).Err()
}

func GetPageCache (rdb *redis.Client) (Cache []model.Page,err error){
	val,err  := rdb.Get(rdbctx,g.PAGE).Result()
	if err != nil {
		return nil,err
	}
	
	if err = json.Unmarshal([]byte(val),&Cache); err != nil {
		return nil,err
	}

	return Cache,err
}

// config
func AddConfigCache(rdb  *redis.Client,config map[string]string) error{
	return rdb.HMSet(rdbctx,g.CONFIG,config).Err()
}

func DelConfigCache(rdb *redis.Client) error{
	return rdb.Del(rdbctx,g.CONFIG).Err()
}

func GetConfigCache(rdb *redis.Client) (Cache map[string]string, err error){	
	return rdb.HGetAll(rdbctx,g.CONFIG).Result()
}