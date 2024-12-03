# Redis
---
数据分冷热，热数据存储在redis（即内存中），大大提高数据访问速度
通常来说呢读数据都是**先从Redis中读取，若没有再去Mysql中读取，后续再写进redis**
<div align=center> <img src=https://notesyzn.oss-cn-hongkong.aliyuncs.com/Redisimage.png width=100%/> </div>


#　Redis持久化
前面说Redis数据在内存中，那么关机不就数据丢失了吗？
<div align=center> <img src=https://notesyzn.oss-cn-hongkong.aliyuncs.com/Redisimage-1.png width=100%/> </div>

Redis在硬盘中存在有`ＡＯＦ和ＲＤＢ文件`,所有写命令都会追加到ＡＯＦ文件中．
－　ＡＯＦ（Append only file）　记录了所有**写**命令
－　ＲＤＢ（Redis Database）　　记录了所有redis数据的二进制信息

>　启动Redis的时候，先会对比AOF RDB看看是不是有什么写操作没有执行，执行后启动Redis　保障了数据的持久化.

## AOF持久化
写操作到日志的持久化方式就是redis中的Append only file 功能.
`AOF`文件一般都是先执行写操作再加载到AOF文件中的，这样有两个好处
- 避免额外的检查开销

AOF两个问题：`写命令没有存入磁盘，服务宕机`和`写命令存入磁盘时间太长导致阻塞下一个写操作`．　这两个问题都和AOF写回时机有关．

### 三种写回策略
 
 <img src=https://notesyzn.oss-cn-hongkong.aliyuncs.com/Redisimage-8.png width=100%/> </div>
>　　Redis 执行完写操作命令后，会将命令追加到 server.aof_buf 缓冲区； 
　　然后通过 write() 系统调用，将 aof_buf 缓冲区的数据写入到 AOF 文件，此时数据并没有写入到硬盘，而是拷贝到了内核缓冲区 page cache，等待内核将数据写入硬盘； 
　　具体内核缓冲区的数据什么时候写入到硬盘，由内核决定。

- `Always`，这个单词的意思是「总是」，所以它的意思是每**次写操作命令执行完后，同步**将 AOF 日志数据写回硬盘； 
- `Everysec`，这个单词的意思是「每秒」，所以它的意思是每次写操作命令执行完后，先将命令写入到 AOF 文件的内核缓冲区，然后**每隔一秒**将缓冲区里的内容写回到硬盘； 
- `No`，意味着不由 Redis 控制写回硬盘的时机，转交给操作系统控制写回的时机，也就是每次写操作命令执行完后，先将命令写入到 AOF 文件的内核缓冲区，再**由操作系统决定何时**将缓冲区内容写回硬盘。

### AOF重写机制
随着写入命令越来越多，AOF文件也就越来越大, 此时重启恢复就会特别的慢，因此Redis提供了重写机制．

AOF 重写机制是在重写时，**读取当前数据库中的所有键值对，然后将每一个键值对用一条命令记录到「新的 AOF 文件」**，等到全部记录完后，就将新的 AOF 文件替换掉现有的 AOF 文件。

重写机制很**耗时**,因此重写操作不能　放在主进程里面，由后台子进程完成．
> 这里使用子进程而不是线程，因为如果是使用线程，多线程之间会共享内存，那么在修改共享内存数据的时候，需要通过加锁来保证数据的安全，而这样就会降低性能。而使用子进程，创建子进程时，父子进程是共享内存数据的，不过这个共享的内存只能以只读的方式，而当父子进程任意一方修改了该共享内存，就会发生「写时复制」，于是父子进程就有了独立的数据副本，就不用加锁来保证数据安全。

还有个问题，重写 AOF 日志过程中，如果`主进程修改了已经存在 key-value`，此时这个 key-value 数据在`子进程的内存数据就跟主进程的内存数据不一致了`，这时要怎么办呢？ 为了解决这种数据不一致问题，Redis 设置了一个 **AOF 重写缓冲区**，这个缓冲区在创建 bgrewriteaof 子进程之后开始使用。 在重写 AOF 期间，当 Redis 执行完一个写命令之后，它会同时将这个写命令写入到 「AOF 缓冲区」和 「AOF 重写缓冲区」。

##　RDB持久化
RDB记录一瞬间的内存数据是二进制文件，而AOF文件记录的是命令操作的日志．因此RDB的恢复效率要比AOg高一些，直接将RDB读入就行．

RDB存在save and bgsave，其中`bgsave是在后台子进程保存快照`.值得一提的是：Redis 的快照是`全量快照`，也就是说每次执行快照，都是把内存中的「`所有数据`」都记录到磁盘中。
> 因此执行快照是一个比较重的操作，如果频率太频繁，可能会对 Redis 性能产生影响。如果频率太低，服务器故障时，丢失的数据会更多。

`bgsave 快照`过程中，如果主线程修改了共享数据，发生了写时复制后，RDB 快照保存的是原本的内存数据，而`主线程刚修改的数据`，是没办法在这一时间写入 RDB 文件的，只能交由下一次的 bgsave 快照。

## AOF　RDB混用
有没有方法不仅有AOF丢失数据少的同时又有RDB恢复速度快的优点呢?
混合使用

*混合使用主要在AOF**重写阶段***
<div align=center> <img src=https://notesyzn.oss-cn-hongkong.aliyuncs.com/Redisimage-9.png width=100%/> </div>

> 当开启了混合持久化时，在 AOF 重写日志时，fork 出来的重写子进程会先将与``主线程共享的内存数据以 RDB 方式写入到 AOF 文件``，然后主线程处理的操作命令会被记录在重写缓冲区里，重写缓冲区里的``增量命令会以 AOF 方式写入到 AOF 文件``，写入完成后通知主进程将新的含有 RDB 格式和 AOF 格式的 AOF 文件替换旧的的 AOF 文件。 也就是说，使用了混合持久化，AOF 文件的前半部分是 RDB 格式的全量数据，后半部分是 AOF 格式的增量数据。
## Redis单线程操作
<img src=https://notesyzn.oss-cn-hongkong.aliyuncs.com/Redisimage-2.png width=100%/> </div>

`代表如果有多个命令等待执行，Redis是排队执行的`
Redis一般为多个Server提供读写服务，为了避免并发问题，Redis被设计为单线程，因此他的操作都是**原子**的

# redis 数据结构详解
##  字符串 (String)
字符串 (String) 是最简单的数据类型，一个键对应一个值。

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
在这个例子中，键 user:1:name 对应值 Alice，这是一对一的关系

### 底层结构
<div align=center> <img src=https://notesyzn.oss-cn-hongkong.aliyuncs.com/Redisimage-3.png width=100%/> </div>

String的基本组成如下，类似于Go中的切片，`len是长度，alloc是容量`
首先sds指针向左获取数据的元信息，知晓数据在buffer中占多大，其次指针向右读取数据

## 哈希 (Hash)

哈希 (Hash) 类似于字典，存储字段和值的映射。它适用于存储对象的多个属性，每个属性对应一个字段和值。
它是为了读取相同属性数据时，减少要Get多个Key的开销设计的．


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
```
> 　同时呢，**Pipeline**也是Hash写入常用的工具．　他的存在是为了在一次redis链接中，执行多条指令，减少网络传输．
> 如果要对Hash数据进行增加操作，不需要全部读出来再操作，`HIncrBy(key, field string, incr int64)`就行
### 底层结构
<div align=center> <img src=https://notesyzn.oss-cn-hongkong.aliyuncs.com/Redisimage-6.png width=100%/> </div>

redis　Hash由哈系表数组组成.

dictht　是常见的Hash，key 通过Hash算法找到嘈位，然后再通过嘈位的链表遍历找到数据．
> Rehash：　当数据太多了，导致嘈位的链表太长，严重影响了查找效率，此时就需要加嘈．　重新创建一dictht，把Ht[0]的数据都转移到Ht[1]中.

>　为了Copy过程中不阻塞用户的访问，那么接需要渐进式移动数据，把Copy操作绑定在用户访问数据操作．

##  列表 (List)

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

### 底层结构
Redis中的ＱuickList结构
<div align=center> <img src=https://notesyzn.oss-cn-hongkong.aliyuncs.com/Redisimage-5.png width=100%/> </div>

> Redis中的List是一个双向链表，同时，Redis为了节省内存空间，在一个节点上存储了多个信息，所有每一个链表节点存储的是 listpack


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

## 有序集合 (ＺSet)
有序集合 (ＺSet) 类似于集合，但每个成员关联一个分数，用于排序。它适用于存储有序的元素集合，如排行榜、时间线等。

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
### 底层结构:
`SKiplist　跳表`
<div align=center> <img src=https://notesyzn.oss-cn-hongkong.aliyuncs.com/Redisimage-7.png width=100%/> </div>

>　跳跃列表的每一层都是一个有序的链表，链表中每个节点都包含两个指针，一个指向同一层的下了一个节点，另一个指向下一层的同一个节点。最低层的链表将包含 zset 中的所有元素。如果说一个元素出现在了某一层，那么低于该层的所有层都将包含这个元素，也就说高层是底层的子集。 

其实就是吧链表分为**许多个子表,Redis中最高四层**，这个子表可以帮助快速访问到目标数据

Zset是　跳表+ Hash　


# Redis　实用举例
## 连续签到
- 业务需求：　每天连续签到，签到数＋１，若没有签到归０
- 用到的知识点： Redis的**Incry**增加１，**ExpireAt**数据过期时间
- 使用到的数据结构:String

## 消息队列
- 业务背景：　当文章通过审核后，如何推送到ＥＳ搜索中使它能被搜索到
- 用到的知识点：　Redis list　消息队列
- 数据结构：list

<div align=center> <img src=https://notesyzn.oss-cn-hongkong.aliyuncs.com/Redisimage-4.png width=100%/> </div>

## Redis实现排行榜
- 业务背景：实现战力排行榜
- 用到数据结构　Ｚset 有序集合
如果使用mysql　读取并排列数据，当访问量过大的时候，很容易崩溃，排列数据非常耗时．


# Go-redis 基本操作
### 1. 连接与配置
- **`NewClient(addr, password string, db int)`**: 创建 Redis 客户端。

### 2. 基本命令
- **`Set(key, value, expiration)`**: 设置键值对。
- **`Get(key)`**: 获取键值。
- **`Del(keys)`**: 删除键。
- **`Exists(keys)`**: 检查键是否存在。
- **`Expire(key, expiration)`**: 设置键的过期时间。
- **`TTL(key)`**: 获取键的剩余生存时间。
- **`Incr(key)`**: 增加键的值。
- **`Decr(key)`**: 减少键的值。

### 3. 字符串操作
- **`Append(key, value)`**: 追加值到键的现有值。
- **`GetRange(key, start, end)`**: 获取键值的子字符串。
- **`SetRange(key, offset, value)`**: 覆盖键值的一部分。
- **`StrLen(key)`**: 获取键值的长度。

### 4. 哈希操作
- **`HSet(key, field, value)`**: 设置哈希表字段和值。
- **`HGet(key, field)`**: 获取哈希表字段值。
- **`HGetAll(key)`**: 获取哈希表所有字段和值。
- **`HDel(key, fields)`**: 删除哈希表字段。
- **`HExists(key, field)`**: 检查哈希表字段是否存在。
- **`HIncrBy(key, field, incr)`**: 增加哈希表字段值。

### 5. 列表操作
- **`LPush(key, values)`**: 插入值到列表头部。
- **`RPush(key, values)`**: 插入值到列表尾部。
- **`LPop(key)`**: 移除并返回列表第一个元素。
- **`RPop(key)`**: 移除并返回列表最后一个元素。
- **`LRange(key, start, stop)`**: 获取列表指定范围的元素。
- **`LLen(key)`**: 获取列表长度。

### 6. 集合操作
- **`SAdd(key, members)`**: 向集合添加成员。
- **`SMembers(key)`**: 获取集合所有成员。
- **`SRem(key, members)`**: 从集合移除成员。
- **`SIsMember(key, member)`**: 检查成员是否存在。
- **`SCard(key)`**: 获取集合成员数量。

### 7. 有序集合操作
- **`ZAdd(key, members)`**: 向有序集合添加成员。
- **`ZRange(key, start, stop)`**: 获取有序集合指定范围的成员。
- **`ZRangeByScore(key, min, max)`**: 根据分数范围获取成员。
- **`ZRem(key, members)`**: 从有序集合移除成员。
- **`ZScore(key, member)`**: 获取有序集合成员的分数。

# Redis　注意事项
## 大key 热key
- 当String类型　值的字节大于10kb就是大key
- 当其他复杂数据结构类型，元素个数大于5000个或总字节数大于10MB就是大Key

`因为Redis是单线程，因此容易导致慢查询，读取成本高，阻塞其他查询操作`，业务侧使用大key会超时报错 

> Ｓtring 如果无法避免使用大key，可以使用拆分kwy，压缩(压缩算法处理ｖ)
> 集合内结构可以考虑，如Zset　区分冷热，只缓存前10页数据，后续数据走db

`热Key `
当一个key值的qbs特别高，几乎是每一个用户都会访问他的时候，那他就是热key

解决方法：　设置localcache，把数据缓存在服务器，不用走redis请求．

## 慢查询
容易导致redis慢查询的操作:
- 批量一次性传入过多的key/value
- zset　复杂度O(log(n)),当大小超过５ｋ，会导致慢查询
- 大key

# 缓存穿透，缓存雪崩

**缓存穿透**：热点数据查询绕过缓存，直接查询数据库
**缓存雪崩**：大量缓存同时过期

`缓存穿透的危害`
- 查询一个**一定不存在**的数据 通常不会缓存不存在的数据，这类查询请求都会直接打到db，如果有系统bug或人为攻击，那么容易导致db响应慢甚至宕机
> 一般redis查询不到的数据查询db

-　缓存过期时 
     在高并发场景下，一个**热key如果过期**，会有大量请求同时击穿至db，容易影响db性能和稳定。
     同一时间有大量key集中过期时，也会导致大量请求落到db上，导致查询变慢，甚至出现db无法响应新的查询

`防止缓存雪崩`
- redis　集群，避免单击单击宕机引起缓存雪崩
- 把缓存时间分散开，这样避免同时多个热key过期，给redis重新建立热key的时间