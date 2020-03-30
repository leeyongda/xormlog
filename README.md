# xormlog
xorm 日志扩展实现(基于logrus日志库)  
https://github.com/sirupsen/logrus

## xorm 版本必须 >= 1.0.0, 其他版本不兼容实现
> ### go get xorm.io/xorm@v1.0.0

### 简单使用教程
```golang
func NewMySQL() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", "dsn")
	if err != nil {
		panic(err)
    }
    logs := logrus.New()
    // 使用自定义日志实现
	logctx := xormlog.NewLogCtx(logs)
    engine.SetLogger(logctx)
    // 需要开启sql输出
    engine.ShowSQL(true)
    
    return engine
}

```
![golang demo](https://carbon.now.sh/?bg=rgba(171%2C184%2C195%2C100)&t=material&wt=none&l=text%2Fx-go&ds=true&dsyoff=20px&dsblur=68px&wc=true&wa=true&pv=48px&ph=32px&ln=false&fl=1&fm=Hack&fs=13px&lh=133%25&si=false&es=2x&wm=false&code=func%2520NewMySQL()%2520*xorm.Engine%2520%257B%250A%2509engine%252C%2520err%2520%253A%253D%2520xorm.NewEngine(%2522mysql%2522%252C%2520%2522dsn%2522)%250A%2509if%2520err%2520!%253D%2520nil%2520%257B%250A%2509%2509panic(err)%250A%2520%2520%2520%2520%257D%250A%2520%2520%2520%2520logs%2520%253A%253D%2520logrus.New()%250A%2520%2520%2520%2520%252F%252F%2520%25E4%25BD%25BF%25E7%2594%25A8%25E8%2587%25AA%25E5%25AE%259A%25E4%25B9%2589%25E6%2597%25A5%25E5%25BF%2597%25E5%25AE%259E%25E7%258E%25B0%250A%2509logctx%2520%253A%253D%2520xormlog.NewLogCtx(logs)%250A%2520%2520%2520%2520engine.SetLogger(logctx)%250A%2520%2520%2520%2520%252F%252F%2520%25E9%259C%2580%25E8%25A6%2581%25E5%25BC%2580%25E5%2590%25AFsql%25E8%25BE%2593%25E5%2587%25BA%250A%2520%2520%2520%2520engine.ShowSQL(true)%250A%2520%2520%2520%2520%250A%2520%2520%2520%2520return%2520engine%250A%257D)