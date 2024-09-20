package model
import "gorm.io/gorm"

// resource 的增删改查
func AddResource (db *gorm.DB,url,name,method string,anonymous bool) (*Resource,error){
	resource := Resource{
		Url: url,
		Name: name,
		Method: method,
		Anonymous: anonymous,
	}

	result := db.Save(&resource)

	return &resource,result.Error
}

func DeleteResource (db *gorm.DB, id int) (int,error){
	result := db.Delete(&Resource{},id)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(db.RowsAffected),result.Error
}

func GetResourceById (db *gorm.DB, id int) (resouce Resource,err error){
	result := db.First(&resouce,id)

	return resouce,result.Error
}

func GetResource (db gorm.DB,url,method string) (resource Resource,err error){
	res := Resource{
		Url: url,
		Method: method,
	}

	result :=db.Where(&res).First(&resource)

	return resource,result.Error
}

// 检查这个resource有没有role使用
func CheckResourceInUse(db *gorm.DB,id int) (bool,error){
	var count int64
	result := db.Model(&RoleResource{}).Where("resource_id = ?",id).Count(&count)
	
	if result.Error != nil{
		return false,result.Error
	}
	return count > 0,nil
}

// 得到特定角色的可用resource
func CheckResourceOfTheRole (db *gorm.DB, rid int) (resource []Resource,err error){
	var role Role
	result := db.Model(&Role{}).Preload("Resources").Take(&role,rid)
	
	return role.Resources,result.Error
}

func CheckResourceHasChild (db *gorm.DB, id int) (bool,error){
	var count int64
	result := db.Model(&Resource{}).Where("parent_id = ?",id).Count(&count)

	return count > 0, result.Error
}

func UpdateResourceAnonymous(db *gorm.DB,id int , anonymous bool) error{
	result := db.Model(&Resource{}).Where("id = ?",id).Update("anonymous",anonymous)
	return result.Error
}

// role
func AddRoleWithResource (db *gorm.DB,name,label string,rs []int) (*Role,error){
	role := Role{
		Name: name,
		Label: label,
	}
	result := db.Create(&role)
	if result.Error != nil{
		return nil,result.Error
	}

	for _,rid := range rs{
		result := db.Model(&RoleResource{}).Create(RoleResource {
			RoleId : role.ID,
			ResourceId : rid,
		})
		if result.Error != nil{
			return nil,result.Error
		}
	}

	return &role,nil
}

func UpdateRoleWithResource(db *gorm.DB,id int,name,label string,rs []int) (*Role,error){
	role := Role{
		Model :Model{ID: id},
		Name: name,
		Label: label,
	}
	
	result := db.Model(&role).Select("name","label").Updates(&role)
	if result.Error != nil {
		return nil,result.Error
	}

	result = db.Delete(&RoleResource{},role.ID)
	if result.Error != nil {
		return nil,result.Error
	}

	for _,rid := range rs{
		result = db.Create(&RoleResource{RoleId: role.ID,ResourceId: rid})
		if result.Error != nil {
			return nil,result.Error
		}
	}

	return &role,nil
}

func DeleteRole(db *gorm.DB,rid int) error {
	result := db.Delete(&RoleMenu{},"role_id = ?",rid)
	if result.Error != nil {
		return result.Error
	}

	result = db.Delete(&RoleResource{},"role_id = ?",rid)
	if result.Error != nil {
		return result.Error
	}

	result = db.Delete(&Role{},rid)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 查找role有没有对应resource的权限或者有没有可以匿名访问的函数
func CheckRoleAuth(db *gorm.DB,url,method string, rid int)(bool,error){
	resource,err := CheckResourceOfTheRole(db,rid)
	if err != nil{
		return false,err
	}

	for _,rs := range resource{
		if(url == rs.Url && method == rs.Method) || rs.Anonymous{
			return true,nil
		}
	}

	return false,nil
}