# go-redis 使用指南

`go-redis` 是一个用于 Go 语言的 Redis 客户端库，它提供了对 Redis 命令的封装，使得在 Go 程序中可以方便地与 Redis 进行交互。`go-redis` 支持 Redis 的所有数据结构和命令。下面是 `go-redis` 中常用的方法和存储的数据结构的简要说明：

## 1. 连接与配置

- **`NewClient(opt *Options) *Client`**: 创建一个新的 Redis 客户端实例。
- **`Options`**: 配置 Redis 连接的选项，包括地址、密码、数据库编号等。

## 2. 基本命令

- **`Set(key string, value interface{}, expiration time.Duration) *StatusCmd`**: 设置键值对，并可以设置过期时间。
- **`Get(key string) *StringCmd`**: 获取指定键的值。
- **`Del(keys ...string) *IntCmd`**: 删除一个或多个键。
- **`Exists(keys ...string) *IntCmd`**: 检查一个或多个键是否存在。
- **`Expire(key string, expiration time.Duration) *BoolCmd`**: 设置键的过期时间。
- **`TTL(key string) *DurationCmd`**: 获取键的剩余生存时间。
- **`Incr(key string) *IntCmd`**: 将键的值增加 1。
- **`Decr(key string) *IntCmd`**: 将键的值减少 1。

## 3. 字符串操作

- **`Append(key string, value string) *IntCmd`**: 将值追加到键的现有值的末尾。
- **`GetRange(key string, start, end int64) *StringCmd`**: 获取键值的子字符串。
- **`SetRange(key string, offset int64, value string) *IntCmd`**: 从指定偏移量开始，覆盖键值的一部分。
- **`StrLen(key string) *IntCmd`**: 获取键值的长度。

## 4. 哈希操作

- **`HSet(key string, values ...interface{}) *IntCmd`**: 设置哈希表中的字段和值。
- **`HGet(key, field string) *StringCmd`**: 获取哈希表中指定字段的值。
- **`HGetAll(key string) *StringStringMapCmd`**: 获取哈希表中所有字段和值。
- **`HDel(key string, fields ...string) *IntCmd`**: 删除哈希表中的一个或多个字段。
- **`HExists(key, field string) *BoolCmd`**: 检查哈希表中是否存在指定字段。
- **`HIncrBy(key, field string, incr int64) *IntCmd`**: 将哈希表中字段的值增加指定的整数。

## 5. 列表操作

- **`LPush(key string, values ...interface{}) *IntCmd`**: 将一个或多个值插入到列表的头部。
- **`RPush(key string, values ...interface{}) *IntCmd`**: 将一个或多个值插入到列表的尾部。
- **`LPop(key string) *StringCmd`**: 移除并返回列表的第一个元素。
- **`RPop(key string) *StringCmd`**: 移除并返回列表的最后一个元素。
- **`LRange(key string, start, stop int64) *StringSliceCmd`**: 获取列表中指定范围的元素。
- **`LLen(key string) *IntCmd`**: 获取列表的长度。

## 6. 集合操作

- **`SAdd(key string, members ...interface{}) *IntCmd`**: 向集合中添加一个或多个成员。
- **`SMembers(key string) *StringSliceCmd`**: 获取集合中的所有成员。
- **`SRem(key string, members ...interface{}) *IntCmd`**: 从集合中移除一个或多个成员。
- **`SIsMember(key string, member interface{}) *BoolCmd`**: 检查成员是否存在于集合中。
- **`SCard(key string) *IntCmd`**: 获取集合的成员数量。

## 7. 有序集合操作

- **`ZAdd(key string, members ...Z) *IntCmd`**: 向有序集合中添加一个或多个成员。
- **`ZRange(key string, start, stop int64) *StringSliceCmd`**: 获取有序集合中指定范围的成员。
- **`ZRangeByScore(key string, opt ZRangeBy) *StringSliceCmd`**: 根据分数范围获取有序集合中的成员。
- **`ZRem(key string, members ...interface{}) *IntCmd`**: 从有序集合中移除一个或多个成员。
- **`ZScore(key, member string) *FloatCmd`**: 获取有序集合中指定成员的分数。

## 8. 发布/订阅

- **`Subscribe(channels ...string) *PubSub`**: 订阅一个或多个频道。
- **`Publish(channel string, message interface{}) *IntCmd`**: 向指定频道发布消息。
- **`PSubscribe(patterns ...string) *PubSub`**: 订阅与指定模式匹配的频道。

## 9. 事务

- **`TxPipeline() Pipeliner`**: 开始一个事务管道。
- **`Exec() ([]Cmder, error)`**: 执行事务中的所有命令。

## 10. 管道

- **`Pipeline() Pipeliner`**: 创建一个管道，允许批量执行多个命令。
- **`Pipelined(fn func(Pipeliner) error) ([]Cmder, error)`**: 在管道中执行一系列命令。

## 11. 脚本

- **`Eval(script string, keys []string, args ...interface{}) *Cmd`**: 执行 Lua 脚本。
- **`EvalSha(sha1 string, keys []string, args ...interface{}) *Cmd`**: 执行缓存的 Lua 脚本。

## 12. 其他

- **`Ping() *StatusCmd`**: 检查与 Redis 服务器的连接是否正常。
- **`Close() error`**: 关闭与 Redis 服务器的连接。

## 存储的数据结构

- **字符串 (String)**: 最简单的数据类型，存储字符串或二进制数据。
- **哈希 (Hash)**: 类似于字典，存储字段和值的映射。
- **列表 (List)**: 有序的字符串列表，支持从两端插入和弹出。
- **集合 (Set)**: 无序的字符串集合，支持集合运算。
- **有序集合 (Sorted Set)**: 类似于集合，但每个成员关联一个分数，用于排序。

## 总结

`go-redis` 提供了丰富的 API 来操作 Redis 的各种数据结构和命令。通过这些方法，你可以轻松地在 Go 程序中实现缓存、队列、分布式锁等功能。

# redis 数据结构详解
##  字符串 (String)
字符串 (String) 是最简单的数据类型，存储字符串或二进制数据。**它是一对一的关系**，即一个键对应一个值。

示例
```go
// 存储字符串
err := rdb.Set(ctx, "user:1:name", "Alice", 0).Err()
if err != nil {
    panic(err)
}

// 获取字符串
name, err := rdb.Get(ctx, "user:1:name").Result()
if err != nil {
    panic(err)
}
fmt.Println("User Name:", name)
```

在这个例子中，键 user:1:name 对应值 Alice，这是一对一的关系。

## 哈希 (Hash)
哈希 (Hash) 类似于字典，存储字段和值的映射。它是**一对多的关系**，即一个键对应多个字段和值。
哈希 (Hash) 类似于字典，存储字段和值的映射。它适用于存储对象的多个属性，每个属性对应一个字段和值。

特点
- 结构化存储: 哈希适合存储结构化的数据，如用户信息、配置信息等。

- 高效读写: 对于单个字段的读写操作非常高效。

- 内存优化: Redis 对哈希的内存使用进行了优化，适合存储大量小对象。
示例

```go
// 存储哈希
err := rdb.HSet(ctx, "user:1", "name", "Alice", "age", 30, "email", "alice@example.com").Err()
if err != nil {
    panic(err)
}

// 获取哈希中的所有字段和值
userInfo, err := rdb.HGetAll(ctx, "user:1").Result()
if err != nil {
    panic(err)
}
fmt.Println("User Info:", userInfo)

```
在这个例子中，键 user:1 对应多个字段和值（name, age, email），这是一对多的关系。

##  列表 (List)
列表 (List) 是有序的字符串列表，支持从两端插入和弹出。它是**一对多的关系**，即一个键对应多个元素。

列表 (List) 是有序的字符串列表，支持从两端插入和弹出。它适用于存储有序的元素集合，如任务队列、消息队列等。

特点
- 有序性: 列表中的元素是有序的，可以按插入顺序访问。

- 双端操作: 支持从列表的两端进行插入和弹出操作，适合实现栈和队列。

- 阻塞操作: 支持阻塞式弹出操作，适合实现消息队列。
- 
示例
```go
// 存储列表
err := rdb.LPush(ctx, "tasks", "task1", "task2", "task3").Err()
if err != nil {
    panic(err)
}

// 获取列表中的所有元素
tasks, err := rdb.LRange(ctx, "tasks", 0, -1).Result()
if err != nil {
    panic(err)
}
fmt.Println("Tasks:", tasks)
```

在这个例子中，键 tasks 对应多个元素（task1, task2, task3），这是一对多的关系。

## 集合 (Set)
集合 (Set) 是无序的字符串集合，支持集合运算。**它是一对多的关系**，即一个键对应多个元素。

集合 (Set) 是无序的字符串集合，支持集合运算。它适用于存储不重复的元素集合，如标签、好友列表等。

特点
- 无序性: 集合中的元素是无序的，不支持按索引访问。

- 唯一性: 集合中的元素是唯一的，不会重复。

- 集合运算: 支持交集、并集、差集等集合运算，适合实现标签、好友关系等功能。
示例

````go
// 存储集合
err := rdb.SAdd(ctx, "tags", "tag1", "tag2", "tag3").Err()
if err != nil {
    panic(err)
}

// 获取集合中的所有元素
tags, err := rdb.SMembers(ctx, "tags").Result()
if err != nil {
    panic(err)
}
fmt.Println("Tags:", tags)

````
在这个例子中，键 tags 对应多个元素（tag1, tag2, tag3），这是一对多的关系。

## 有序集合 (Sorted Set)
有序集合 (Sorted Set) 类似于集合，但每个成员关联一个分数，用于排序。它是一对多的关系，即一个键对应多个成员和分数。

有序集合 (Sorted Set) 类似于集合，但每个成员关联一个分数，用于排序。它适用于存储有序的元素集合，如排行榜、时间线等。

特点
- 有序性: 有序集合中的元素是按分数排序的，可以按分数范围访问。

- 唯一性: 有序集合中的元素是唯一的，不会重复。

- 分数排序: 每个元素关联一个分数，适合实现排行榜、时间线等功能。
示例

```go
// 存储有序集合
err := rdb.ZAdd(ctx, "scores", &redis.Z{Score: 90, Member: "Alice"}, &redis.Z{Score: 85, Member: "Bob"}).Err()
if err != nil {
    panic(err)
}

// 获取有序集合中的所有成员和分数
scores, err := rdb.ZRangeWithScores(ctx, "scores", 0, -1).Result()
if err != nil {
    panic(err)
}
fmt.Println("Scores:", scores)

```
在这个例子中，键 scores 对应多个成员和分数（Alice: 90, Bob: 85），这是一对多的关系。

## 总结
字符串 (String): 一对一的关系，一个键对应一个值。

哈希 (Hash): 一对多的关系，一个键对应多个字段和值。

列表 (List): 一对多的关系，一个键对应多个元素。

集合 (Set): 一对多的关系，一个键对应多个元素。

有序集合 (Sorted Set): 一对多的关系，一个键对应多个成员和分数。

通过这些数据结构，Redis 提供了丰富的功能来满足不同的应用需求。


# go-redis 数据结构的隐式选择
在 Redis 中，数据结构的选择是隐式的，它取决于你使用的命令。Redis 不会要求你在存储数据之前显式地指定数据结构，而是根据你执行的命令自动选择合适的数据结构。

- 当你使用 SET 命令时，Redis 会自动将数据存储为字符串（String）。

- 当你使用 HSET 命令时，Redis 会自动将数据存储为哈希（Hash）。

- 当你使用 LPUSH 或 RPUSH 命令时，Redis 会自动将数据存储为列表（List）。

- 当你使用 SADD 命令时，Redis 会自动将数据存储为集合（Set）。

- 当你使用 ZADD 命令时，Redis 会自动将数据存储为有序集合（Sorted Set）。

包括但是我限于这几个方法.只要使用数据结构对应的这些方法，就能确定数据的存储结构。