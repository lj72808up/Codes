#### 数据类型
1. 字符串  
 可以存储数字，字符串。本质是一个字节数组，因为支持2种数据类型，所以有2种类型的操作。 
    * INCR，DECR等（数字类型）
    * getRange，setRange，setBit，getBit（字符串型）
    
2. 列表   
  因为列表是有序序列，所以支持
    * 左右两端 push 操作
    * 左右两端 pop 操作
    * 按照 index 访问： lIndex 
    * 范围操作： lRange，lTrim
    
   阻塞队列相关，pop 操作可以阻塞一段时间，等待有元素到来在返回
    * bLPop,bRPop
    
   从 sorce 列表弹出后推入 dist 列表的事务操作
    * rPopLPush
    * bRpopLpush
    
3. 集合 set
   无序不重复
    * sAdd，sRem 增删集合中的成员
    * sMembers 集合所有成员
    * sIsMember 是否是集合中的成员
    * sMove 将 src 中的成员 item 移除添加到 dist 成员中
    
   集合间的操作
    * sDiff，sInter，sUnion 求差集，交集，并集
    * sDiffStore，sInterStore，sUnionStore 求差集，交集，并集后存储
    
4. Hash 散列    
   单个操作
    * hmGet，hmSet，hDel 获取，添加，删除
    * hLen 返回长度
    * hIncrBy 键对应的值增加 （hIncrByFloat）
    
   批量操作
    * hExists
    * hKeys，hVals 所有的 keys，values
    
5. 有序集合 zSet   
   对键操作
    * zAdd,zRem 增删
    * zCard 成员数量
    
   分值操作
    * zIncrBy 给成员的分值增加
    * zCount 分值介于min，max之间
    * zScore 返回成员分值
    * zRangeByScore 返回分值介于min，max之间的成员
    * zRemRangeByScore 删除分值介于min，max之间的成员
    
   排名操作
    * zRank 返回成员的排名
    * zRevRank 返回成员从大到小排列后，成员的排名
    * zRevRange 成员从大到小排列后，返回排名在 min，max 之间的成员
    * zRange 返回排名在 min，max 之间的成员
    
   集合之间的操作
    * zInterStore 默认的交集函数取 sum
    * zUnionStore 聚合函数可以制定 min

#### 持久化
1. 基于快照的 bgSave   
  适用于即使丢失一部分数据，也不会造成问题的应用 
    * bgSave采用 fork 子进程，子进程创建完毕后开始导出快照，父进程继续响应客户端命令
    * save 命令不 fork 子进程， 自己到处快照文件
    * 在 redis 占用几十 GB 内存后， bgSave 创建的子进程因为也要占用几十 GB，导致刀量使用操作系统虚拟内存，导出会非常慢。  
      这种情况推荐在夜间使用 save 命令，暂时不响应所有客户端请求，知道导出快照完毕