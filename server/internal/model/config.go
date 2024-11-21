package model

import (
	"strconv"

	"gorm.io/gorm"
)

type Config struct {
	Model
	Key			string			`gorm:"unique;type:varchar(256)" json:"key"`
	Value		string			`gorm:"type:varchar(256)"  json:"value"`
	Desc		string			`gorm:"type:varchar(256)"  json:"desc"`
}

func GetConfigMap(db *gorm.DB) (map[string]string,error){
	var list []Config
	result := db.Find(&list)
	if result.Error != nil {
		return nil,result.Error
	}

	configMap := make(map[string]string)
	for _,config := range list{
		configMap[config.Key] = config.Value
	}

	return configMap,nil
}

func UpdateConfigMap(db *gorm.DB,configMap map[string]string) error{
	for key,value := range configMap{
		result := db.Model(&Config{}).Where("key",key).Update("value",value)
		
		if result.Error != nil{
			return result.Error
		}
	}
	return nil
}

func FindOrCreateConfig(db *gorm.DB,key,value string) error {
	var config Config
	result := db.Where("key",key).FirstOrCreate(&config)
	if result.Error != nil {
		return result.Error
	}

	
	config.Value = value
	result = db.Save(&config)

	if result.Error != nil {
		return result.Error
	}
	
	return nil
}


func GetValueByKey(db *gorm.DB,key string) (string){
	var config Config
	result := db.Model(&config).Where("key",key).First(&config)
	if result.Error != nil {
		return ""
	}

	return config.Value
}

func GetConfigBool(db *gorm.DB,key string) bool{
	value := GetValueByKey(db,key)
	if value == ""{
		return false
	}
	return value == "true"
}

func GetConfigInt(db *gorm.DB,key string) int{
	value:= GetValueByKey(db,key)
	if value == ""{
		return 0
	}

	result,err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return result 
}
