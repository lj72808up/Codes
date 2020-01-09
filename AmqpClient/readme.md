1. 通过配置文件配置多个数据库
    (1) 配置文件中有1个datasource name的列表
    (2) 对每个列表配置一个druid数据源
2. 程序启动时, 创建一个单利对象: 持有一个"datasource name"和"datasource"的映射
3.使用PrintWriter, BufferedOutputStream 输出文件


4. 明确bufferedWriter在复制时产生的数组复制
(1) 如果需要写入的string的长度 > 内部buffer的长度:
        a. flushbuffer() 用内部的fileoutputstream把缓冲的buf全都写出去, 写出的长度用一个count来记录当前内部buffer的长度
        b. 目标string对应的byte[]数组直接用底层的fileoutputstream全部写出
        c. 写完返回
(2) 目标string对应的bytes[]数组长度 > 内部buffer剩余(buffersize-count)的长度
        a. flushbuffer() 用内部的fileoutputstream把缓冲的buf全都写出去
        b. 写出后, 当前的buffer内部已经空了

(3) 讲目标string的bytes[]复制到内部的缓冲buffer(System.arraycopy)
(4) count += string对应bytes[]的长度


### TODO
rabbit mq的消费者和生产者的ack机制