// page数据库存储的是每一个网址的背景图片

package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

type Page struct{}

// 获取page和对应的背景图片
func (*Page) GetList(c *gin.Context){
	//先从redis缓存中获取 更快
	data,err := GetPageCache(GetRDB(c))
	if data != nil && err == nil{
		ReturnSuccess(c,data)
		return
	}
	
	switch err {
	case redis.Nil:
		break
	default: 
		ReturnError(c,g.ErrRedisOp,err)
		return
	}	
	
	// redis失败 从db获取后加入缓存
	data,_,err = model.GetPageList(GetDB(c))
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	if err := AddPageCache(GetRDB(c), data); err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c,data)
}

// 新增或更新背景
//操作完成后都需要清空缓存
func (*Page) SaveAndUpdate(c *gin.Context){
	var page model.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	data,err := model.SaveOrCreatePage(GetDB(c),page.ID,page.Name,page.Label,page.Cover)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	// 删除缓存
	if err = DelPageCache(GetRDB(c)); err != nil{
		ReturnError(c,g.ErrRedisOp,err)
		return
	}
	
	ReturnSuccess(c,data)	
}

// 删除背景
func (*Page) Delete(c *gin.Context){
	var ids []int
	if err := c.ShouldBind(&ids); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	result := GetDB(c).Delete(&model.Page{},"id IN ?",ids)
	if result.Error!= nil {
		ReturnError(c,g.ErrDbOp,result.Error)
		return
	}
	// 删除缓存
	if err := DelPageCache(GetRDB(c)); err!= nil{
		ReturnError(c,g.ErrRedisOp,err)
		return
	}

	ReturnSuccess(c,nil)
}