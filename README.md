# Hello !
这是一个用Go开发的简单博客服务器,用于学习Go的web开发。
开发技术：Gin + Gorm

### 目录说明：
cmd :入口函数所在目录

configs：配置文件所在目录，需要自己创建该目录。
该目录下存放config.dev.json。
dev可修改为prod、test对应不同的环境。
默认环境为dev，可通过设置操作系统的变量来修改环境，或者在main.go中修改环境。

internal：
代码主要存放的目录。

| 目录名        | 作用          | 是否可用 | 上次修改 |  
|------------|-------------|------|------|
| config     | 配置对象所在目录    | yes  | --   |
| database   | 数据库模块       | yes  | --   |
| handler    | Gin 处理器所在目录 | yes  | --   |
| middleware | 中间件所在目录     | yes  | --   |
| model      | 数据模型所在目录    | yes  | --   |
| router     | 路由所在目录      | yes  | --   |
|service|其它服务所在目录| no   | --   |

pkg目录:

util:中使用的主要工具

|方法名| 作用                                       | 是否可用 |上次修改|
|----|------------------------------------------|-----|----|
|AppendLogBuff| 用于存放写入数据库的日志消息，目前用于错误处理中间件，避免一直写数据库      |
|FlushLogBuffer| 把AppendLyesogBuff中的消息存入数据库 ，持久化记录，防止内存爆炸 |
|GenerateToken| 生成Token ,用于权限验证中间件                       | yes |
|GetLimiter| 获取访问速率限制器，用于LimitMiddleware 中间件，防止请求超速   | yes |

文件说明：
Server.js 系统暴露的接口文件，前端使用。
config.dev.json 系统配置文件，修改对应的配置项目。

使用：
git xxx
构建->发布

最后：
系统存在BUG,会持续完善和更新。喜欢的朋友可以点个小星星。






