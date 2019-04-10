# 日志相关

## 1. 标准输出、标准出错日志
尽量不要在标准输出、标准出错输出非panic级别的日志，将标准输出、标准出错留给panic关键日志，方便跟踪进程状态。
```
fmt.Println("xxx")                             // 标准输出
log.Println("xxx")                             // 标准出错
var logger = github.com/lovego/config.Logger() // 标准出错
logger.Info("xxx")
```

只要写到标准输出、标准出错的日志，都是写到了docker容器的标准输出、标准出错的，
因此可以使用`docker logs`或`xiaomei app logs`查看。

## 2. 输出到文件的日志
```
var logger = github.com/lovego/config.NewLogger("xxx.log")
logger.Info("xxx")
```
输出到文件的日志，如果收集到了ElasticSearch，使用Kibana查看。
否则使用`xiaomei app shell`进入docker容器，日志文件在log目录下。

## 3. 跟HTTP请求绑定的日志
```
router.Get(`/`, func(c *goa.Context) {                    
  ctx := c.Context()
  tracer.Log(ctx, `xxx`)           // 日志
  tracer.Tag(ctx, `key`, `value`)  // 标签
  c.Json(map[string]string{`hello`: config.DeployName()}) 
})                                                        
```
跟请求绑定的日志，默认已经收集到了ElasticSearch，因此使用Kibana查看。

