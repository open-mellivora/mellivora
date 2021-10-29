# mini_spider

## 程序设计

![](document/mini_spider.png)

1. spider 产生请求发送给 engine
2. engine 把请求发给 scheduler 调度
3. engine 从 scheduler 取到待执行的请求
4. engine 把请求发给 middleware 进行包装
5. middleware 把请求发给 downloader 下载
6. downloader 下载
7. middleware 对 response 进行处理
8. middleware 把 response 发送给 engine
9. engine 把返回给 spider,回到开头

## TODO

- Context 序列化, 支持远程任务
- 优化 Spider 接口
- 增加扩展功能
- 添加更多常用中间件
