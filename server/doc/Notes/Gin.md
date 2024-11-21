# http方法 POST PUT PATCH 的区别
---
## POST
## 用途:
     用于向服务器提交数据，通常用于创建新资源。

## 场景:

- 创建新资源: 当你需要向服务器发送数据并希望服务器创建一个新的资源时，使用 POST 方法。例如，提交表单数据、上传文件、创建新用户等。

- 批量操作: 当你需要对多个资源进行批量操作时，也可以使用 POST 方法。例如，批量上传文件、批量创建用户等。

## 特点:

- 非幂等性: POST 请求是非幂等的，即多次调用可能会产生不同的结果。例如，每次 POST 请求都可能创建一个新的资源。

- 请求体: 请求参数通过请求体传递。

- 无缓存: 请求不会被缓存。

- 无书签: 请求不能被书签保存。


## PUT
### 用途
用于更新服务器上的资源，通常用于完全替换现有资源。

### 场景:

- 完全替换资源: 当你需要更新服务器上的现有资源，并且希望用新的数据完全替换旧的数据时，使用 PUT 方法。例如，更新用户信息、修改文章内容等。

- 创建资源: 如果资源不存在，PUT 方法也可以用于创建资源，但通常不推荐这种用法。

### 特点:

- 幂等性: PUT 请求是幂等的，即多次调用会产生相同的结果。例如，多次调用 PUT 请求更新同一个资源，结果是相同的。

- 请求体: 请求参数通过请求体传递。

- 无缓存: 请求不会被缓存。

- 无书签: 请求不能被书签保存。

##  PATCH
## 用途: 
用于对资源进行部分更新。

## 场景:

部分更新资源: 当你只需要更新资源的部分内容时，使用 PATCH 方法。例如，只更新用户的某个字段，而不需要替换整个资源。

## 特点:

- 非幂等性: PATCH 请求通常是非幂等的，即多次调用可能会产生不同的结果。例如，每次 PATCH 请求都可能更新资源的不同部分。

- 请求体: 请求参数通过请求体传递。

- 无缓存: 请求不会被缓存。

- 无书签: 请求不能被书签保存。


## 总结
POST: 用于创建新资源或批量操作，非幂等。

PUT: 用于完全替换现有资源，幂等。

PATCH: 用于部分更新资源，通常非幂等。

## handler 处理 GET 和 POST 请求中数据的区别
对于 **GET** 请求，参数通常通过 URL 传递。Gin 提供了 `c.Query` 和 `c.Param` 方法来解析这些参数。

- c.Query：用于解析 URL 中的查询参数（即 ?key=value 部分）。 对多个form数据合适

- c.Param：用于解析 URL 路径中的参数（即路由中的占位符部分）。 

对于 **POST**请求，参数可以通过多种方式传递，如表单参数、JSON 参数等。Gin 提供了相应的方法来解析这些参数。

-`c.PostForm`：用于解析表单参数（即 application/x-www-form-urlencoded 或 multipart/form-data 类型的请求体）。

-`c.BindJSON`：用于解析 JSON 参数（即 application/json 类型的请求体）。


# Gin 使用细节
## 动态参数
例如在 `"/comment/like/:comment_id"` 中，就把like后面的参数作为动态参数，在handler中使用 `c.Param("comment_id")` 来获取。