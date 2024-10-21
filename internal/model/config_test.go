package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigMap(t *testing.T) {
	db := setup(t)

	configs := []Config{
		{Key: "name", Value: "Blog", Desc: "姓名"},
		{Key: "age", Value: "12", Desc: "年龄"},
		{Key: "enabled", Value: "true", Desc: "是否可用"},
	}
	db.Create(&configs)

	data, err := GetConfigMap(db)
	assert.Nil(t, err)
	assert.Len(t, data, 3)
	assert.Equal(t, "Blog", data["name"])
	assert.Equal(t, "12", data["age"])
	assert.Equal(t, "true", data["enabled"])
}

func TestUpdateConfigMap(t *testing.T) {
	db := setup(t)
	configs := []Config{
		{Key: "name", Value: "Blog", Desc: "姓名"},
		{Key: "age", Value: "12", Desc: "年龄"},
		{Key: "enabled", Value: "true", Desc: "是否可用"},
	}
	db.Create(&configs)

	m := map[string]string{
		"name":    "Alice",
		"age":     "15",
		"enabled": "false",
		"dump":    "dump", // 无效数据
	}
	err := UpdateConfigMap(db, m)
	assert.Nil(t, err)

	data, err := GetConfigMap(db)
	assert.Nil(t, err)
	assert.Len(t, data, 3)
	assert.Equal(t, "Alice", data["name"])
	assert.Equal(t, "15", data["age"])
	assert.Equal(t, "false", data["enabled"])
}

func TestConfigSetGet(t *testing.T) {
	db := setup(t)

	FindOrCreateConfig(db, "name", "AAA")

	val:= GetValueByKey(db, "name")
	assert.Equal(t, "AAA", val)

	m, _ := GetConfigMap(db)
	assert.Len(t, m, 1)
}

func TestCheckConfig(t *testing.T) {
	db := setup(t)

	{
		FindOrCreateConfig(db, "name", "AAA")
		val:= GetValueByKey(db, "name")
		assert.Equal(t, "AAA", val)
	}

	{
		FindOrCreateConfig(db, "name", "BBB")
		val:= GetValueByKey(db, "name")
		assert.Equal(t, "BBB", val)
	}

}
