## Limiter
一组简单的限流模型

### 常见模型
  - 令牌桶 Token Bucket done
  - 漏桶 Leakey Bucket done
  - 固定窗口计数 todo
  - 滑动窗口计数 todo
  
### Token Bucket
初始化令牌桶实例:
```go
tb := NewTokenBucket(&TbConfig{
    QPS:    10,
    MaxCap: 20,
})
```
`QPS: 为速率最小单位`

`MaxCap: 桶最大容量, 保证突发大流量冲击`

领取令牌:
```go
if err := tb.Take(); err != nil {
    atomic.AddInt64(&unpass, 1)
} else {
    atomic.AddInt64(&pass, 1)
}
```
`err == nil 为领取到令牌`

### Leaky Bucket
初始化漏桶实例:
```go
lb := NewLeakyBucket(&LbConfig{
    Rate:     1,
    MaxSlack: 0,
})
```
`Rate: 为速率，每秒请求数`

`MaxSlack: 为最大松弛度，默认为0不开启`


### 更多信息
![avatar](https://github.com/xiaoxuz/fql/blob/master/wechat.jpg)

  

  
