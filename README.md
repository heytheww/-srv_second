# Srv_second
这是一个简单商品秒杀微服务项目，力求尽量使用纯Go语言标准库开发，并努力保持代码整洁、支持高并发的特性。目前已经完成的功能：
- [x] 登录，签发token
- [x] 限流
- [x] 使用 redis 记录秒杀成功用户名单
- [ ] 使用 redis消息队列 和 mysql 异步生成订单信息
- [ ] 错误处理
- [ ] 单元测试

目前已经实现的功能是：  
1 秒杀开始前，服务读取goodsId_count为0，直接返回未开始；秒杀开始时，将goodsId_count改为 >0 ，标志秒杀开始，goodsId_count为 0 后，拦截所有请求。  
2 goodsId_count保存在redis中，实现外部系统修改goodsId_count马上触发秒杀开始。  
3 秒杀成功，下单用户ID放入set中。


## Background
这个项目很早就构思了，当时是基于Node.js去实现，整体代码比较啰嗦并且高并发的支持也不是很强。现在使用 Go 标准库实现，天然支持超高并发，系统资源占用小。  
注意，jwt的部分不是本项目的重点，仅仅是为了项目完整性，即 登录--秒杀--记录结果。jwt 签发的 token 将会交给我的另一个开源项目：heytheww/
srv_gateway 业务网关 去完成鉴权。


## Install
下载源码，准备好 Go 1.18 环境

## Usage
在项目根目录：
```
$ go run .
```

## Contributing
提个Issues说出你的需求 或 提个PR实现你的想法，我会测试并整合到本项目。

## License
MIT © heytheww