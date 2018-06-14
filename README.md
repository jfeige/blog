# blog

基于golang写的一款简单的个人博客。框架采用Gin，数据库采用mysql,采用redis作为缓存,文本编辑器采用的是Editmd,整个工程包含了前台和后台。

### 配置
```
./conf/blog.ini 是应用和mysql，redis的一些基本配置
./conf/blog-log.xml 是log4go需要的配置文件
```
工程根目录下，有个env.sh文件，赋给可执行权限后，直接执行即可:
```
./env.sh
```
### 前台地址:

http://101.201.253.64:8090/   

### 后台地址

http://101.201.253.64:8090/login

帐号:admin  密码:123456



项目中使用到的第三方包:

框架: gopkg.in/gin-gonic/gin.v1

redis: github.com/garyburd/redigo/redis

mysql: github.com/go-sql-driver/mysql

日志: github.com/alecthomas/log4go

配置文件: github.com/jfeige/lconfig


