# xormlog
xorm 日志扩展实现(基于logrus日志库)  
https://github.com/sirupsen/logrus

## xorm 版本必须 >= 1.0.0, 其他版本不兼容实现
> ### go get xorm.io/xorm@v1.0.0

### 简单使用教程
<details><summary>main.go</summary>
<p>

```golang
func NewMySQL() *xorm.Engine {
  engine, err := xorm.NewEngine("mysql", "dsn")
  if err != nil {
    panic(err)
  }
  logs := logrus.New()
  // 使用自定义日志实现
  logctx := xormlog.NewLogCtx(logs)
  // 默认同步输出sql日志，也可以开启异步输出sql日志
  // logctx.AsyncShowSQL(true)
  // defer logctx.Close()
  engine.SetLogger(logctx)
    // 需要开启sql输出
  engine.ShowSQL(true)
  return engine
}

``` 
</p>
</details>