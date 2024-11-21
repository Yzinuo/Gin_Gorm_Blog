# 声明模型

GORM 通过将 Go 结构体（Go structs） **映射到数据库表**来简化数据库交互。

- CreatedAt UPdatedAt 两个字段是特殊字段，当记录被创建或更新时，GORM 会自动向内填充当前时间。
- **主键**：GORM 使用一个名为ID 的字段作为每个模型的默认主键。

- **表名**：默认情况下，GORM 将结构体名称转换为 snake_case 并为表名加上复数形式。 例如，一个 User 结构体在数据库中的表名变为 users 。

- **列名**：GORM 自动将结构体字段名称转换为 snake_case 作为数据库中的列名。

---

# 查询
## 单个查询（First Last）
First and Last 方法会按主键排序找到第一条记录和最后一条记录 (分别)。 只有在目标**struct 是指针或者通过 db.Model() 指定 model 时，该方法才有效**。且没有找到记录时，它会返回 ErrRecordNotFound 错误

### 根据主键查询
单个查询时，可以传入第二个参数：主键值。

```go
var user User
db.First(&user, 1) // 查询主键为 1 的记录
```

当目标对象有一个主键值时，将使用主键构建查询条件，例如：
```go
var user = User{ID: 10}
db.First(&user)
// SELECT * FROM users WHERE id = 10;

var result User
db.Model(User{ID: 10}).First(&result)
```

## 条件查询
**GORM 的查询构建器会收集所有的查询条件，并在你调用 Find、First、Take、Scan 等方法时才会真正执行查询Where。**

> 如果对象设置了主键，条件查询将不会覆盖主键的值，而是用 And 连接条件。
> 写gorm就规范一些，结构体就只指定主键值，剩余的条件都用Where声明
```go
var user = User{ID: 10}
db.Where("id = ?", 20).First(&user)
// SELECT * FROM users WHERE id = 10 and id = 20 ORDER BY id ASC LIMIT 1

```
### struct && map 条件查询（Where方法）
GORM **只会将结构体对象中的非零值作为查询条件**。 要使零值也作为查询条件，需要使用 map 作为查询条件。

### 内联查询
不使用Where Model，直接使用First Find 等方法。
```go
db.Find(&user, "name = ?", "jinzhu")
// SELECT * FROM users WHERE name = "jinzhu";
```

### Joins 添加关联表
```go
db.Joins("Company").Find(&users)
// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` LEFT JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id`;
```

# 更新
## Save
Save 会保存所有的字段，即使字段是零值
如果保存值没有主键，它将执行Create，否则它将执行Update。
```go
    user := User{ID:1,Name: "jinzhu", Age: 18}
    db.Save(&user)
```
**不要将 Save 和 Model一同使用**, 这是 未定义的行为。

## Updates 批量更新
如果你**执行一个没有任何条件**的批量更新，GORM 默认不会运行，并且会返回 ErrMissingWhereClause 错误
```go
db.Model(&User{}).Update("name", "jinzhu").Error // gorm.ErrMissingWhereClause
```

# 删除
Delete的方法使用和First一样
- 前面model指定表并且查询完成后 调用delete批量删除
- delete 传入 struct，指定表的同时指定要删除的记录的主键
- delte第一个参数指定表，第二个参数（未声明的情况）是目标记录的主键
- 内联查询  第二个参数是查询条件。
删除一条记录时，删除对象需要指定主键，**否则会触发 批量删除**

当你试图执行不带任何条件的批量删除时，GORM将不会运行并返回ErrMissingWhereClause 错误。 如果硬是要这么做，用原生sql

# 像First， Find， Delte这些可以自己指定表的方法，都可以不需要Model

# Belongs to 定义外键，子表
belongs to 会与另一个模型建立了一对一的连接（子表逻辑上来看）。
```go
// 子表的model中要包含外键，和父表model
type User struct {
  gorm.Model
  Name      string
  CompanyID int
  Company   Company
}

type Company struct {
  ID   int
  Name string
}
```
Blong to 和 Has one， Has three查询主体不同。
例如一个用户有一张信用卡。
以信用卡为查询主体，就要使用Belong to模式 卡属于用户。  以用户为查询主体，就要使用Has one模式，用户拥有一张信用卡。
两者在实现上正好是互补的 连贯的。

# Has Many 一对多
```go
// User 有多张 CreditCard，UserID 是外键
type User struct {
  gorm.Model
  CreditCards []CreditCard  `gorm:"foreignKey:UserID"`  
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint
}
```
gorm:"foreignKey:UserID" 作用是定义子表中的外键，即哪一个字段和我们的主键关联。

# ManytoMany
```go
// User 拥有并属于多种 language，`user_languages` 是连接表
type User struct {
  gorm.Model
  Languages []*Language `gorm:"many2many:user_languages;"`
}

type Language struct {
  gorm.Model
  Name string
  Users []*User `gorm:"many2many:user_languages;"`
}
```

# Scopes
作用域允许你复用通用的逻辑，这种共享逻辑需要定义为类型func(*gorm.DB) *gorm.DB
```go
func Paginate(page, size int) func(db *gorm.DB) *gorm.DB{
	return func (db *gorm.DB) *gorm.DB{
		if page <=0{
			page = 1
		}
		switch {
		case  size >= 100:
			size = 100	
		case  size <= 10 : 
			size = 10
		}
		offset := (page - 1) *size
		return db.Offset(offset).Limit(size)
	}
}
```

# 不可寻指错误
---
在 Go 语言中，只有**可寻址的值**（addressable value）才能被赋值。不可寻址的值包括：

- 字面量（如 123、"hello"）

- 常量

- 函数的返回值（除非返回的是指针）

- 映射（map）中的元素

- 切片（slice）中的元素（除非切片本身是可寻址的）

我们需要确保在gorm操作中，Find，First，Scan等方法的参数是可寻址的，这样gorm反射操作才不会出问题。
而`指针是可寻址的`，因此Find，First，Scan等方法的参数是指针，而不是变量本身。