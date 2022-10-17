### 二、查询缓存

查询缓存 (Query Cache)：以 key-value 形式存在内存。key -> SQL 查询语句，value -> SQL 查询语句结果。

只要表有更新操作，查询缓存就会清空。
关闭缓存：query_cache_type=DEMAND
MySQL 8.0 移除了 Server 层的查询缓存，不是 InnoDB 存储引擎的 buffer pool.

### 三、解析 SQL

#### 解析器
