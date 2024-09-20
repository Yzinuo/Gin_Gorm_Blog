package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	db := setup(t)

	// 添加资源
	res1,_ := AddResource(db,"/api/method1","Api1","GET",false)
	res2,_ :=  AddResource(db,"/api/method2","Api2","PUT",false)

	// 创造有这两个资源的角色
	role ,err := AddRoleWithResource(db,"admin","管理员",[]int{res1.ID,res2.ID})
	assert.Nil(t,err)
	assert.Equal(t,role.Label,"管理员")
	
	rs,_ := GetResourceIdsByRoleId(db,role.ID)
	assert.Equal(t,len(rs),2)

	// 更新这个角色
	role,err = UpdateRoleWithResource(db,role.ID,"admin","超级用户",[]int{res1.ID})
	assert.Nil(t,err)
	rs,_ = GetResourceIdsByRoleId(db,role.ID)
	assert.Equal(t,len(rs),1)

	// 鉴权
	flag,_ := CheckRoleAuth(db,"/api/method1","GET",role.ID)
	assert.Equal(t,flag,true)
	flag,_ = CheckRoleAuth(db,"/api/method2","PUT",role.ID)
	assert.Equal(t,flag,false)
}